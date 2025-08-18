package tokenizer

import (
	"fmt"
	"strings"

	"github.com/pkoukk/tiktoken-go"
)

// TiktokenTokenizer wraps the tiktoken library for OpenAI and Anthropic providers
type TiktokenTokenizer struct {
	encoding *tiktoken.Tiktoken
	provider string
}

// NewTiktokenTokenizer creates a new tokenizer for OpenAI and Anthropic providers
func NewTiktokenTokenizer(provider string) (*TiktokenTokenizer, error) {
	encoding := getEncodingForProvider(provider)

	enc, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		return nil, fmt.Errorf("failed to create tiktoken tokenizer for provider %s: %w", provider, err)
	}

	return &TiktokenTokenizer{
		encoding: enc,
		provider: provider,
	}, nil
}

// CountTokens counts the number of tokens in the given text
func (t *TiktokenTokenizer) CountTokens(text string) (int, error) {
	tokens := t.encoding.Encode(text, nil, nil)
	return len(tokens), nil
}

// GetModelName returns the provider name (for compatibility)
func (t *TiktokenTokenizer) GetModelName() string {
	return t.provider
}

// getEncodingForProvider determines the appropriate encoding for a given provider
func getEncodingForProvider(provider string) string {
	providerLower := strings.ToLower(provider)

	switch providerLower {
	case "anthropic", "openai":
		// Both Anthropic and OpenAI use cl100k_base encoding
		return "cl100k_base"
	default:
		// Default to cl100k_base
		return "cl100k_base"
	}
}

// IsTiktokenProvider checks if the provider should use tiktoken
func IsTiktokenProvider(provider string) bool {
	providerLower := strings.ToLower(provider)
	return providerLower == "anthropic" || providerLower == "openai"
}

// SupportedTiktokenProviders lists providers that use tiktoken
var SupportedTiktokenProviders = []string{
	"anthropic",
	"openai",
}
