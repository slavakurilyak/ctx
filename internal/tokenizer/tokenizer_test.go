package tokenizer

import (
	"strings"
	"sync"
	"testing"
)

// MockTokenizer is a mock implementation for testing
type MockTokenizer struct {
	CountTokensFunc  func(string) (int, error)
	GetModelNameFunc func() string
}

func (m *MockTokenizer) CountTokens(text string) (int, error) {
	if m.CountTokensFunc != nil {
		return m.CountTokensFunc(text)
	}
	// Default: simple word count
	return len(strings.Fields(text)), nil
}

func (m *MockTokenizer) GetModelName() string {
	if m.GetModelNameFunc != nil {
		return m.GetModelNameFunc()
	}
	return "mock-model"
}

// MockTokenizerFactory creates mock tokenizers for testing
type MockTokenizerFactory struct {
	CreateTokenizerFunc func(string) (Tokenizer, error)
}

func (m *MockTokenizerFactory) CreateTokenizer(modelName string) (Tokenizer, error) {
	if m.CreateTokenizerFunc != nil {
		return m.CreateTokenizerFunc(modelName)
	}
	return &MockTokenizer{}, nil
}

func TestTokenizerInterface(t *testing.T) {
	tests := []struct {
		name      string
		tokenizer Tokenizer
		text      string
		wantCount int
		wantModel string
	}{
		{
			name: "mock tokenizer",
			tokenizer: &MockTokenizer{
				CountTokensFunc: func(text string) (int, error) {
					return 42, nil
				},
				GetModelNameFunc: func() string {
					return "test-model"
				},
			},
			text:      "Hello, world!",
			wantCount: 42,
			wantModel: "test-model",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count, err := tt.tokenizer.CountTokens(tt.text)
			if err != nil {
				t.Fatalf("CountTokens() error = %v", err)
			}
			if count != tt.wantCount {
				t.Errorf("CountTokens() = %v, want %v", count, tt.wantCount)
			}

			model := tt.tokenizer.GetModelName()
			if model != tt.wantModel {
				t.Errorf("GetModelName() = %v, want %v", model, tt.wantModel)
			}
		})
	}
}

func TestTokenizerCache(t *testing.T) {
	factory := &MockTokenizerFactory{
		CreateTokenizerFunc: func(modelName string) (Tokenizer, error) {
			return &MockTokenizer{
				GetModelNameFunc: func() string {
					return modelName
				},
			}, nil
		},
	}

	cache := NewTokenizerCache(factory)

	// Test creating and caching
	tok1, err := cache.GetOrCreate("model1")
	if err != nil {
		t.Fatalf("GetOrCreate() error = %v", err)
	}

	// Should return cached instance
	tok2, err := cache.GetOrCreate("model1")
	if err != nil {
		t.Fatalf("GetOrCreate() error = %v", err)
	}

	// Compare pointers to ensure same instance
	if tok1 != tok2 {
		t.Error("GetOrCreate() should return cached instance")
	}

	// Test different model
	tok3, err := cache.GetOrCreate("model2")
	if err != nil {
		t.Fatalf("GetOrCreate() error = %v", err)
	}

	if tok3 == tok1 {
		t.Error("GetOrCreate() should create new instance for different model")
	}

	// Test cache size
	if cache.Size() != 2 {
		t.Errorf("Size() = %v, want 2", cache.Size())
	}

	// Test Has
	if !cache.Has("model1") {
		t.Error("Has() should return true for cached model")
	}
	if cache.Has("model3") {
		t.Error("Has() should return false for uncached model")
	}

	// Test Clear
	cache.Clear()
	if cache.Size() != 0 {
		t.Errorf("Size() after Clear() = %v, want 0", cache.Size())
	}
}

func TestTokenizerCacheConcurrency(t *testing.T) {
	factory := &MockTokenizerFactory{}
	cache := NewTokenizerCache(factory)

	// Test concurrent access
	var wg sync.WaitGroup
	models := []string{"model1", "model2", "model3"}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			modelName := models[index%len(models)]
			_, err := cache.GetOrCreate(modelName)
			if err != nil {
				t.Errorf("Concurrent GetOrCreate() error = %v", err)
			}
		}(i)
	}

	wg.Wait()

	// Should have exactly 3 models cached
	if cache.Size() != 3 {
		t.Errorf("Size() after concurrent access = %v, want 3", cache.Size())
	}
}

func TestIsGeminiModel(t *testing.T) {
	tests := []struct {
		provider string
		want     bool
	}{
		{"gemini", true},
		{"Gemini", true}, // Case insensitive
		{"GEMINI", true}, // Case insensitive
		{"openai", false},
		{"anthropic", false},
		{"claude", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			if got := IsGeminiModel(tt.provider); got != tt.want {
				t.Errorf("IsGeminiModel(%q) = %v, want %v", tt.provider, got, tt.want)
			}
		})
	}
}

// TestIsTiktokenProvider tests if a provider uses tiktoken for tokenization
func TestIsTiktokenProvider(t *testing.T) {
	tests := []struct {
		provider string
		want     bool
	}{
		{"openai", true},
		{"anthropic", true},
		{"gemini", false},
		{"unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			// Check by trying to create tokenizer and seeing if it's tiktoken-based
			factory := &DefaultTokenizerFactory{}
			tokenizer, err := factory.CreateTokenizer(tt.provider)

			if tt.provider == "unknown" {
				if err == nil {
					t.Errorf("Expected error for unknown provider %q", tt.provider)
				}
				return
			}

			if err != nil && tt.want {
				t.Errorf("Failed to create tokenizer for %q: %v", tt.provider, err)
				return
			}

			_, isTiktoken := tokenizer.(*TiktokenTokenizer)
			if isTiktoken != tt.want {
				t.Errorf("Provider %q: got tiktoken=%v, want %v", tt.provider, isTiktoken, tt.want)
			}
		})
	}
}

func TestDefaultTokenizerFactory(t *testing.T) {
	factory := &DefaultTokenizerFactory{}

	tests := []struct {
		provider  string
		wantError bool
	}{
		{"openai", false},
		{"anthropic", false},
		{"gemini", false},
		{"unknown", true}, // Unknown provider should error
		{"", true},        // Empty provider should error
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			_, err := factory.CreateTokenizer(tt.provider)
			if (err != nil) != tt.wantError {
				t.Errorf("CreateTokenizer(%q) error = %v, wantError %v", tt.provider, err, tt.wantError)
			}
		})
	}
}

func TestGetSupportedProviders(t *testing.T) {
	providers := GetSupportedProviders()
	if len(providers) == 0 {
		t.Error("GetSupportedProviders() should return non-empty list")
	}

	// Check that all expected providers are included
	expected := map[string]bool{
		"anthropic": false,
		"openai":    false,
		"gemini":    false,
	}

	for _, provider := range providers {
		if _, ok := expected[provider]; ok {
			expected[provider] = true
		}
	}

	for provider, found := range expected {
		if !found {
			t.Errorf("GetSupportedProviders() should include %q", provider)
		}
	}
}

func TestIsProviderSupported(t *testing.T) {
	tests := []struct {
		provider string
		want     bool
	}{
		{"openai", true},
		{"OpenAI", true}, // Case insensitive
		{"anthropic", true},
		{"ANTHROPIC", true}, // Case insensitive
		{"gemini", true},
		{"Gemini", true}, // Case insensitive
		{"unknown-provider", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.provider, func(t *testing.T) {
			if got := IsProviderSupported(tt.provider); got != tt.want {
				t.Errorf("IsProviderSupported(%q) = %v, want %v", tt.provider, got, tt.want)
			}
		})
	}
}
