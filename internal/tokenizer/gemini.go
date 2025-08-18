package tokenizer

import (
	"fmt"
	"strings"

	"cloud.google.com/go/vertexai/genai"
	"cloud.google.com/go/vertexai/genai/tokenizer"
)

// GeminiTokenizer wraps Google's Vertex AI tokenizer for Gemini provider
type GeminiTokenizer struct {
	tok      *tokenizer.Tokenizer
	provider string
}

// NewGeminiTokenizer creates a new tokenizer for Gemini provider
func NewGeminiTokenizer(provider string) (*GeminiTokenizer, error) {
	// Use a default Gemini model for tokenization
	defaultModel := "gemini-1.5-pro"

	tok, err := tokenizer.New(defaultModel)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini tokenizer for provider %s: %w", provider, err)
	}

	return &GeminiTokenizer{
		tok:      tok,
		provider: provider,
	}, nil
}

// CountTokens counts the number of tokens in the given text
func (g *GeminiTokenizer) CountTokens(text string) (int, error) {
	resp, err := g.tok.CountTokens(genai.Text(text))
	if err != nil {
		return 0, fmt.Errorf("failed to count tokens: %w", err)
	}
	return int(resp.TotalTokens), nil
}

// GetModelName returns the provider name (for compatibility)
func (g *GeminiTokenizer) GetModelName() string {
	return g.provider
}

// IsGeminiModel checks if the provider is Gemini (kept for backward compatibility)
func IsGeminiModel(provider string) bool {
	providerLower := strings.ToLower(provider)
	return providerLower == "gemini"
}

// SupportedGeminiProvider is the Gemini provider
var SupportedGeminiProvider = "gemini"
