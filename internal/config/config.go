package config

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// LimitsConfig holds all resource-limit configurations.
// Pointers are used to distinguish between a zero value and an unset value.
type LimitsConfig struct {
	MaxTokens         *int64 `yaml:"max_tokens,omitempty"`
	MaxOutputBytes    *int64 `yaml:"max_output_bytes,omitempty"`
	MaxLines          *int64 `yaml:"max_lines,omitempty"`
	MaxPipelineStages *int   `yaml:"max_pipeline_stages,omitempty"`
}

type Config struct {
	TokenModel        string
	DefaultTimeout    time.Duration
	OutputFormat      string
	PrettyOutput      bool
	CacheDir          string
	MaxTokens         int64 // This will be deprecated in favor of Limits.MaxTokens
	NoTokens          bool
	NoHistory         bool
	NoTelemetry       bool
	NoTelemetrySource string // New field to track the source
	Limits            LimitsConfig
	Auth              *AuthConfig         `yaml:"auth,omitempty"`
	Installation      *InstallationConfig `yaml:"installation,omitempty"`
}

// InstallationConfig tracks how ctx was installed and update preferences
type InstallationConfig struct {
	Method              string        `yaml:"method,omitempty"`                // "install-script", "go-install", "manual", "pre-built"
	AutoUpdateCheck     bool          `yaml:"auto_update_check,omitempty"`     // Whether to check for updates automatically
	LastUpdateCheck     time.Time     `yaml:"last_update_check,omitempty"`     // Last time we checked for updates
	UpdateCheckInterval time.Duration `yaml:"update_check_interval,omitempty"` // How often to check (default: 24h)
	SkipVersions        []string      `yaml:"skip_versions,omitempty"`         // Versions to skip
}

type AuthConfig struct {
	APIKey      string `yaml:"api_key,omitempty"` // Note: This will just be a placeholder; actual key is in keychain
	APIEndpoint string `yaml:"api_endpoint,omitempty"`
	Tier        string `yaml:"tier,omitempty"`
	ExpiresAt   string `yaml:"expires_at,omitempty"`
}

// Define source constants
const (
	SourceDefault = "default"
	SourceEnv     = "environment variable"
	SourceFlag    = "command-line flag"
)

var defaultConfig = &Config{
	TokenModel:        "anthropic",
	DefaultTimeout:    2 * time.Minute,
	OutputFormat:      "json",
	CacheDir:          "",
	NoTelemetrySource: SourceDefault,
	Auth: &AuthConfig{
		// APIEndpoint is loaded from environment or config file
		APIEndpoint: "",
	},
}

func Get() *Config {
	cfg := &Config{
		TokenModel:   getEnvOrDefault("CTX_TOKEN_MODEL", defaultConfig.TokenModel),
		OutputFormat: getEnvOrDefault("CTX_OUTPUT_FORMAT", defaultConfig.OutputFormat),
		CacheDir:     getCacheDir(),
	}

	// Parse timeout
	if timeoutStr := os.Getenv("CTX_TIMEOUT"); timeoutStr != "" {
		if timeout, err := time.ParseDuration(timeoutStr); err == nil {
			cfg.DefaultTimeout = timeout
		} else {
			cfg.DefaultTimeout = defaultConfig.DefaultTimeout
		}
	} else {
		cfg.DefaultTimeout = defaultConfig.DefaultTimeout
	}

	// Parse CTX_PRETTY
	if prettyStr := os.Getenv("CTX_PRETTY"); prettyStr == "true" {
		cfg.PrettyOutput = true
	}

	return cfg
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getCacheDir() string {
	// Check for tiktoken cache dir
	if dir := os.Getenv("TIKTOKEN_CACHE_DIR"); dir != "" {
		return dir
	}

	// Use default cache directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(homeDir, ".cache", "tiktoken")
}

// LoadConfigFromFile loads configuration from ~/.config/ctx/config.yaml
func LoadConfigFromFile() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(homeDir, ".config", "ctx", "config.yaml")

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil // No config file, return nil
	}

	// Read the file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Parse YAML
	var fileConfig struct {
		TokenModel     string       `yaml:"token_model,omitempty"`
		DefaultTimeout string       `yaml:"default_timeout,omitempty"`
		OutputFormat   string       `yaml:"output_format,omitempty"`
		CacheDir       string       `yaml:"cache_dir,omitempty"`
		NoTokens       bool         `yaml:"no_tokens,omitempty"`
		NoHistory      bool         `yaml:"no_history,omitempty"`
		NoTelemetry    bool         `yaml:"no_telemetry,omitempty"`
		Limits         LimitsConfig `yaml:"limits,omitempty"`
		Auth           *AuthConfig  `yaml:"auth,omitempty"`
	}

	if err := yaml.Unmarshal(data, &fileConfig); err != nil {
		return nil, err
	}

	// Build config from file
	cfg := &Config{}

	if fileConfig.TokenModel != "" {
		cfg.TokenModel = fileConfig.TokenModel
	}
	if fileConfig.OutputFormat != "" {
		cfg.OutputFormat = fileConfig.OutputFormat
	}
	if fileConfig.CacheDir != "" {
		cfg.CacheDir = fileConfig.CacheDir
	}
	if fileConfig.DefaultTimeout != "" {
		if d, err := time.ParseDuration(fileConfig.DefaultTimeout); err == nil {
			cfg.DefaultTimeout = d
		}
	}

	cfg.NoTokens = fileConfig.NoTokens
	cfg.NoHistory = fileConfig.NoHistory
	cfg.NoTelemetry = fileConfig.NoTelemetry
	cfg.Limits = fileConfig.Limits
	cfg.Auth = fileConfig.Auth

	return cfg, nil
}

// NewFromFlagsAndEnv creates a Config object by layering flag, environment, and default values.
func NewFromFlagsAndEnv(cmd *cobra.Command) *Config {
	// 1. Start with defaults
	cfg := Get()                          // Assumes Get() provides default values
	cfg.NoTelemetrySource = SourceDefault // Initialize with default source

	// 2. Layer on file configuration
	if fileConfig, err := LoadConfigFromFile(); err == nil && fileConfig != nil {
		// Merge file config into current config
		if fileConfig.TokenModel != "" {
			cfg.TokenModel = fileConfig.TokenModel
		}
		if fileConfig.OutputFormat != "" {
			cfg.OutputFormat = fileConfig.OutputFormat
		}
		if fileConfig.CacheDir != "" {
			cfg.CacheDir = fileConfig.CacheDir
		}
		if fileConfig.DefaultTimeout != 0 {
			cfg.DefaultTimeout = fileConfig.DefaultTimeout
		}
		if fileConfig.NoTokens {
			cfg.NoTokens = fileConfig.NoTokens
		}
		if fileConfig.NoHistory {
			cfg.NoHistory = fileConfig.NoHistory
		}
		if fileConfig.NoTelemetry {
			cfg.NoTelemetry = fileConfig.NoTelemetry
			cfg.NoTelemetrySource = "config file"
		}

		// Merge limits
		cfg.Limits = fileConfig.Limits

		// Also set the deprecated MaxTokens if provided in Limits
		if fileConfig.Limits.MaxTokens != nil {
			cfg.MaxTokens = *fileConfig.Limits.MaxTokens
		}

		// Merge Auth configuration
		if fileConfig.Auth != nil {
			cfg.Auth = fileConfig.Auth
		} else {
			// Initialize empty Auth config if not in file
			cfg.Auth = &AuthConfig{}
		}

		// Merge Installation configuration
		if fileConfig.Installation != nil {
			cfg.Installation = fileConfig.Installation
		} else {
			// Initialize default Installation config if not in file
			cfg.Installation = &InstallationConfig{
				Method:              "unknown",
				AutoUpdateCheck:     false,
				UpdateCheckInterval: 24 * time.Hour,
			}
		}
	}

	// 3. Layer on Environment Variables
	if val := os.Getenv("CTX_TOKEN_MODEL"); val != "" {
		cfg.TokenModel = val
	}
	if val := os.Getenv("CTX_TIMEOUT"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			cfg.DefaultTimeout = d
		}
	}
	isPrivateEnv := os.Getenv("CTX_PRIVATE") == "true"
	noHistoryEnv := os.Getenv("CTX_NO_HISTORY") == "true"
	noTelemetryEnv := os.Getenv("CTX_NO_TELEMETRY") == "true"
	noTokensEnv := os.Getenv("CTX_NO_TOKENS") == "true"

	// Check environment variables first for telemetry
	if noTelemetryEnv || isPrivateEnv {
		cfg.NoTelemetry = true
		// Be specific about the source
		if isPrivateEnv {
			cfg.NoTelemetrySource = "CTX_PRIVATE " + SourceEnv
		} else {
			cfg.NoTelemetrySource = "CTX_NO_TELEMETRY " + SourceEnv
		}
	}

	// Handle other env-based configs
	cfg.NoHistory = isPrivateEnv || noHistoryEnv
	cfg.NoTokens = noTokensEnv

	// Handle max tokens from environment
	if maxTokensStr := os.Getenv("CTX_MAX_TOKENS"); maxTokensStr != "" {
		if mt, err := strconv.ParseInt(maxTokensStr, 10, 64); err == nil && mt > 0 {
			cfg.MaxTokens = mt
			// Also update the Limits struct
			cfg.Limits.MaxTokens = &mt
		}
	}

	// Handle other limit environment variables
	if maxOutputBytesStr := os.Getenv("CTX_MAX_OUTPUT_BYTES"); maxOutputBytesStr != "" {
		if mob, err := strconv.ParseInt(maxOutputBytesStr, 10, 64); err == nil && mob > 0 {
			cfg.Limits.MaxOutputBytes = &mob
		}
	}

	if maxLinesStr := os.Getenv("CTX_MAX_LINES"); maxLinesStr != "" {
		if ml, err := strconv.ParseInt(maxLinesStr, 10, 64); err == nil && ml > 0 {
			cfg.Limits.MaxLines = &ml
		}
	}

	if maxPipelineStagesStr := os.Getenv("CTX_MAX_PIPELINE_STAGES"); maxPipelineStagesStr != "" {
		if mps, err := strconv.Atoi(maxPipelineStagesStr); err == nil && mps > 0 {
			cfg.Limits.MaxPipelineStages = &mps
		}
	}

	// Handle API endpoint from environment (overrides file config)
	if apiEndpoint := os.Getenv("CTX_API_ENDPOINT"); apiEndpoint != "" {
		if cfg.Auth == nil {
			cfg.Auth = &AuthConfig{}
		}
		cfg.Auth.APIEndpoint = apiEndpoint
	}

	// Handle CTX_PRETTY environment variable (but only if flag not explicitly set)
	if prettyStr := os.Getenv("CTX_PRETTY"); prettyStr == "true" && !cmd.Flags().Changed("pretty") {
		cfg.PrettyOutput = true
	}

	// 3. Layer on Flags (highest precedence)
	if cmd.Flags().Changed("token-model") {
		cfg.TokenModel, _ = cmd.Flags().GetString("token-model")
	}
	if cmd.Flags().Changed("timeout") {
		cfg.DefaultTimeout, _ = cmd.Flags().GetDuration("timeout")
	}
	if cmd.Flags().Changed("output") {
		cfg.OutputFormat, _ = cmd.Flags().GetString("output")
	}
	if cmd.Flags().Changed("pretty") {
		cfg.PrettyOutput, _ = cmd.Flags().GetBool("pretty")
	}
	if cmd.Flags().Changed("max-tokens") {
		mt, _ := cmd.Flags().GetInt64("max-tokens")
		cfg.MaxTokens = mt
		// Also update the Limits struct
		if mt > 0 {
			cfg.Limits.MaxTokens = &mt
		}
	}

	if cmd.Flags().Changed("max-output-bytes") {
		mob, _ := cmd.Flags().GetInt64("max-output-bytes")
		if mob > 0 {
			cfg.Limits.MaxOutputBytes = &mob
		}
	}

	if cmd.Flags().Changed("max-lines") {
		ml, _ := cmd.Flags().GetInt64("max-lines")
		if ml > 0 {
			cfg.Limits.MaxLines = &ml
		}
	}

	if cmd.Flags().Changed("max-pipeline-stages") {
		mps, _ := cmd.Flags().GetInt("max-pipeline-stages")
		if mps > 0 {
			cfg.Limits.MaxPipelineStages = &mps
		}
	}

	// Handle boolean flags with proper source tracking
	if cmd.Flags().Changed("private") {
		isPrivate, _ := cmd.Flags().GetBool("private")
		if isPrivate {
			cfg.NoHistory = true
			cfg.NoTelemetry = true
			cfg.NoTelemetrySource = "--private " + SourceFlag
		}
	}

	if cmd.Flags().Changed("no-history") {
		if v, _ := cmd.Flags().GetBool("no-history"); v {
			cfg.NoHistory = true
		}
	}

	if cmd.Flags().Changed("no-telemetry") {
		if v, _ := cmd.Flags().GetBool("no-telemetry"); v {
			cfg.NoTelemetry = true
			cfg.NoTelemetrySource = "--no-telemetry " + SourceFlag
		}
	}

	if cmd.Flags().Changed("no-tokens") {
		if v, _ := cmd.Flags().GetBool("no-tokens"); v {
			cfg.NoTokens = true
		}
	}

	// If telemetry is still enabled at this point, ensure the source reflects that
	if !cfg.NoTelemetry {
		cfg.NoTelemetrySource = SourceDefault
	}

	return cfg
}

// getConfigPath returns the path to the config file
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "ctx", "config.yaml"), nil
}

// SaveConfig saves the configuration to the config file
func (c *Config) SaveConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	// Convert to file config struct (for YAML serialization)
	fileConfig := struct {
		TokenModel   string              `yaml:"token_model,omitempty"`
		Timeout      string              `yaml:"timeout,omitempty"`
		OutputFormat string              `yaml:"output_format,omitempty"`
		PrettyOutput bool                `yaml:"pretty_output,omitempty"`
		NoTokens     bool                `yaml:"no_tokens,omitempty"`
		NoHistory    bool                `yaml:"no_history,omitempty"`
		NoTelemetry  bool                `yaml:"no_telemetry,omitempty"`
		Limits       LimitsConfig        `yaml:"limits,omitempty"`
		Auth         *AuthConfig         `yaml:"auth,omitempty"`
		Installation *InstallationConfig `yaml:"installation,omitempty"`
	}{
		TokenModel:   c.TokenModel,
		OutputFormat: c.OutputFormat,
		PrettyOutput: c.PrettyOutput,
		NoTokens:     c.NoTokens,
		NoHistory:    c.NoHistory,
		NoTelemetry:  c.NoTelemetry,
		Limits:       c.Limits,
		Auth:         c.Auth,
		Installation: c.Installation,
	}

	if c.DefaultTimeout > 0 {
		fileConfig.Timeout = c.DefaultTimeout.String()
	}

	data, err := yaml.Marshal(fileConfig)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// SetInstallationMethod sets the installation method and saves the config
func (c *Config) SetInstallationMethod(method string) error {
	if c.Installation == nil {
		c.Installation = &InstallationConfig{}
	}
	c.Installation.Method = method

	// Set reasonable defaults for auto-update
	if method == "install-script" || method == "manual" {
		c.Installation.AutoUpdateCheck = true
	} else {
		c.Installation.AutoUpdateCheck = false
	}

	if c.Installation.UpdateCheckInterval == 0 {
		c.Installation.UpdateCheckInterval = 24 * time.Hour
	}

	return c.SaveConfig()
}
