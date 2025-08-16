package auth

import (
	"github.com/zalando/go-keyring"
)

const serviceName = "ctx-cli"

// KeychainManager handles secure credential storage
type KeychainManager struct{}

// NewKeychainManager creates a new keychain manager
func NewKeychainManager() *KeychainManager {
	return &KeychainManager{}
}

// StoreAPIKey stores an API key in the system keychain
func (k *KeychainManager) StoreAPIKey(username, apiKey string) error {
	return keyring.Set(serviceName, username, apiKey)
}

// GetAPIKey retrieves an API key from the system keychain
func (k *KeychainManager) GetAPIKey(username string) (string, error) {
	return keyring.Get(serviceName, username)
}

// ClearAPIKey removes an API key from the system keychain
func (k *KeychainManager) ClearAPIKey(username string) error {
	return keyring.Delete(serviceName, username)
}