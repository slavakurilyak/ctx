package enricher

import (
	"os"
	"strings"
	"sync"

	"github.com/pkoukk/tiktoken-go"
)

var (
	tokenCounter     *TokenCounter
	tokenCounterOnce sync.Once
	tokenCounterErr  error
)

type TokenCounter struct {
	encoding *tiktoken.Tiktoken
	model    string
}

func CountTokens(text string) (*int, error) {
	tokenCounterOnce.Do(func() {
		tokenCounter, tokenCounterErr = initTokenCounter()
	})

	if tokenCounterErr != nil {
		return nil, tokenCounterErr
	}

	if tokenCounter == nil {
		return nil, nil
	}

	count := tokenCounter.Count(text)
	return &count, nil
}

func initTokenCounter() (*TokenCounter, error) {
	model := getTokenModel()

	// Map model names to their encodings
	encoding := getEncodingForModel(model)

	enc, err := tiktoken.GetEncoding(encoding)
	if err != nil {
		// Try to get encoding for the specific model
		enc, err = tiktoken.EncodingForModel(model)
		if err != nil {
			// If all fails, return nil (token counting disabled)
			return nil, nil
		}
	}

	return &TokenCounter{
		encoding: enc,
		model:    model,
	}, nil
}

func (tc *TokenCounter) Count(text string) int {
	tokens := tc.encoding.Encode(text, nil, nil)
	return len(tokens)
}

func getTokenModel() string {
	model := os.Getenv("CTX_TOKEN_MODEL")
	if model == "" {
		// Default to Claude 4.1 Opus as specified in the README
		model = "claude-4.1-opus"
	}
	return model
}

func getEncodingForModel(model string) string {
	modelLower := strings.ToLower(model)

	// Claude models use cl100k_base (same as GPT-4)
	if strings.Contains(modelLower, "claude") {
		return "cl100k_base"
	}

	// GPT-4 and GPT-3.5-turbo use cl100k_base
	if strings.Contains(modelLower, "gpt-4") ||
		strings.Contains(modelLower, "gpt-3.5-turbo") ||
		strings.Contains(modelLower, "gpt-35-turbo") {
		return "cl100k_base"
	}

	// GPT-3 models use p50k_base
	if strings.Contains(modelLower, "davinci") ||
		strings.Contains(modelLower, "curie") ||
		strings.Contains(modelLower, "babbage") ||
		strings.Contains(modelLower, "ada") {
		return "p50k_base"
	}

	// Default to cl100k_base (most modern encoding)
	return "cl100k_base"
}
