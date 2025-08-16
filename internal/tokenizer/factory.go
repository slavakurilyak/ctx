package tokenizer

import (
	"fmt"
	"os"
	"strings"
)

// DefaultTokenizerFactory is the default implementation of TokenizerFactory
type DefaultTokenizerFactory struct{}

// CreateTokenizer creates the appropriate tokenizer based on the provider name
func (f *DefaultTokenizerFactory) CreateTokenizer(provider string) (Tokenizer, error) {
	if provider == "" {
		return nil, fmt.Errorf("provider name cannot be empty")
	}
	
	// Normalize provider name
	providerLower := strings.ToLower(provider)
	
	switch providerLower {
	case "anthropic", "openai":
		// Both use tiktoken with cl100k_base encoding
		return NewTiktokenTokenizer(provider)
	case "gemini":
		return NewGeminiTokenizer(provider)
	default:
		return nil, fmt.Errorf("unsupported provider: %s (supported: anthropic, openai, gemini)", provider)
	}
}

// NewTokenizer is a convenience function that creates a tokenizer using the default factory
func NewTokenizer(provider string) (Tokenizer, error) {
	factory := &DefaultTokenizerFactory{}
	return factory.CreateTokenizer(provider)
}

// NewTokenizerFromEnv creates a tokenizer based on the CTX_TOKEN_MODEL environment variable
func NewTokenizerFromEnv() (Tokenizer, error) {
	provider := os.Getenv("CTX_TOKEN_MODEL")
	if provider == "" {
		// Default to Anthropic provider
		provider = "anthropic"
	}
	return NewTokenizer(provider)
}

// GetSupportedProviders returns a list of all supported providers
func GetSupportedProviders() []string {
	return []string{"anthropic", "openai", "gemini"}
}

// IsProviderSupported checks if a provider is supported
func IsProviderSupported(provider string) bool {
	providerLower := strings.ToLower(provider)
	for _, supported := range GetSupportedProviders() {
		if supported == providerLower {
			return true
		}
	}
	return false
}