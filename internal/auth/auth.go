package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/slavakurilyak/ctx/internal/config"
	"gopkg.in/yaml.v3"
)

// Manager handles authentication operations
type Manager struct {
	config   *config.Config
	apiKey   string
	keychain *KeychainManager
}

// NewManager creates a new authentication manager
func NewManager(cfg *config.Config) *Manager {
	return &Manager{
		config:   cfg,
		keychain: NewKeychainManager(),
	}
}

// Login authenticates a user with the given API key
func (m *Manager) Login(apiKey string) error {
	// Validate API key with backend
	userInfo, err := m.validateAPIKey(apiKey)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	// Store API key in keychain
	if err := m.keychain.StoreAPIKey(userInfo.Email, apiKey); err != nil {
		return fmt.Errorf("failed to store credentials: %w", err)
	}

	// Update config file with subscription details
	if err := m.updateConfigFile(userInfo); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	m.apiKey = apiKey
	return nil
}

// Logout clears the user's credentials
func (m *Manager) Logout() error {
	// Get current user email from config
	if m.config.Auth == nil || m.config.Auth.APIKey == "" {
		return fmt.Errorf("not logged in")
	}

	// Clear from keychain
	email := m.getCurrentUserEmail()
	if email != "" {
		if err := m.keychain.ClearAPIKey(email); err != nil {
			// Continue even if keychain clear fails
			fmt.Fprintf(os.Stderr, "Warning: could not clear keychain: %v\n", err)
		}
	}

	// Clear from config file
	if err := m.clearConfigFile(); err != nil {
		return fmt.Errorf("failed to clear config: %w", err)
	}

	m.apiKey = ""
	return nil
}

// GetAPIKey retrieves the stored API key
func (m *Manager) GetAPIKey() (string, error) {
	if m.apiKey != "" {
		return m.apiKey, nil
	}

	// Try to get from keychain
	email := m.getCurrentUserEmail()
	if email == "" {
		return "", fmt.Errorf("not logged in")
	}

	apiKey, err := m.keychain.GetAPIKey(email)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve credentials: %w", err)
	}

	m.apiKey = apiKey
	return apiKey, nil
}

// IsAuthenticated checks if the user is authenticated
func (m *Manager) IsAuthenticated() bool {
	_, err := m.GetAPIKey()
	return err == nil
}

// GetAccountInfo retrieves current account information
func (m *Manager) GetAccountInfo() (*UserInfo, error) {
	if m.config.Auth == nil {
		return nil, fmt.Errorf("not logged in")
	}

	// Validate current API key
	apiKey, err := m.GetAPIKey()
	if err != nil {
		return nil, err
	}

	return m.validateAPIKey(apiKey)
}

// UserInfo represents user account information
type UserInfo struct {
	Email      string       `json:"email"`
	Tier       string       `json:"tier"`
	ExpiresAt  time.Time    `json:"expires_at"`
	ValidUntil string       `json:"valid_until"`
	Pricing    *PricingInfo `json:"pricing,omitempty"`
}

// PricingInfo represents pricing information
type PricingInfo struct {
	Plan         string  `json:"plan"`           // "individual" or "team"
	Currency     string  `json:"currency"`       // USD, EUR, etc.
	Amount       float64 `json:"amount"`         // Current price
	BillingCycle string  `json:"billing_cycle"`  // monthly, annual
	Seats        int     `json:"seats,omitempty"` // Number of seats for team plans
	NextBilling  string  `json:"next_billing,omitempty"`
}

// validateAPIKey validates an API key with the backend
func (m *Manager) validateAPIKey(apiKey string) (*UserInfo, error) {
	endpoint := ""
	if m.config != nil && m.config.Auth != nil && m.config.Auth.APIEndpoint != "" {
		endpoint = m.config.Auth.APIEndpoint
	}
	
	// API endpoint must be configured via environment or config
	if endpoint == "" {
		// Try to get from environment as fallback
		endpoint = os.Getenv("CTX_API_ENDPOINT")
		if endpoint == "" {
			return nil, fmt.Errorf("API endpoint not configured. Set CTX_API_ENDPOINT environment variable or configure in ~/.config/ctx/config.yaml")
		}
	}

	req, err := http.NewRequest("GET", endpoint+"/v1/validate", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid API key")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("validation failed with status: %d", resp.StatusCode)
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// updateConfigFile updates the config file with user information
func (m *Manager) updateConfigFile(userInfo *UserInfo) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".config", "ctx", "config.yaml")
	
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return err
	}

	// Read existing config or create new
	var configData map[string]interface{}
	if data, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(data, &configData); err != nil {
			return err
		}
	} else {
		configData = make(map[string]interface{})
	}

	// Get existing auth config to preserve api_endpoint
	existingAuth, _ := configData["auth"].(map[string]interface{})
	apiEndpoint := ""
	if existingAuth != nil && existingAuth["api_endpoint"] != nil {
		apiEndpoint = existingAuth["api_endpoint"].(string)
	}
	
	// If no existing endpoint, try to get from environment
	if apiEndpoint == "" {
		apiEndpoint = os.Getenv("CTX_API_ENDPOINT")
	}
	
	// Update auth section
	configData["auth"] = map[string]interface{}{
		"api_key":      "***", // Placeholder - actual key is in keychain
		"tier":         userInfo.Tier,
		"expires_at":   userInfo.ValidUntil,
		"email":        userInfo.Email,
		"api_endpoint": apiEndpoint,
	}

	// Write back to file
	data, err := yaml.Marshal(configData)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// clearConfigFile removes auth information from the config file
func (m *Manager) clearConfigFile() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(homeDir, ".config", "ctx", "config.yaml")
	
	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		// No config file, nothing to clear
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var configData map[string]interface{}
	if err := yaml.Unmarshal(data, &configData); err != nil {
		return err
	}

	// Preserve api_endpoint but clear auth credentials
	if authData, ok := configData["auth"].(map[string]interface{}); ok {
		// Keep only api_endpoint if it exists
		if endpoint := authData["api_endpoint"]; endpoint != nil {
			configData["auth"] = map[string]interface{}{
				"api_endpoint": endpoint,
			}
		} else {
			delete(configData, "auth")
		}
	}

	// Write back to file
	data, err = yaml.Marshal(configData)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// getCurrentUserEmail retrieves the current user's email from config
func (m *Manager) getCurrentUserEmail() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	configPath := filepath.Join(homeDir, ".config", "ctx", "config.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return ""
	}

	var configData map[string]interface{}
	if err := yaml.Unmarshal(data, &configData); err != nil {
		return ""
	}

	authData, ok := configData["auth"].(map[string]interface{})
	if !ok {
		return ""
	}

	email, ok := authData["email"].(string)
	if !ok {
		return ""
	}

	return email
}