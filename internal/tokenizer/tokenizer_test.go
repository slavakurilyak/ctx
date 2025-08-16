package tokenizer

import (
	"strings"
	"sync"
	"testing"
)

// MockTokenizer is a mock implementation for testing
type MockTokenizer struct {
	CountTokensFunc func(string) (int, error)
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
		model string
		want  bool
	}{
		{"gemini-1.5-flash", true},
		{"gemini-1.5-pro", true},
		{"gemini-2.0-flash-exp", true},
		{"Gemini-1.5-Flash", true}, // Case insensitive
		{"gpt-4", false},
		{"claude-3-opus", false},
		{"", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			if got := IsGeminiModel(tt.model); got != tt.want {
				t.Errorf("IsGeminiModel(%q) = %v, want %v", tt.model, got, tt.want)
			}
		})
	}
}

func TestIsTiktokenModel(t *testing.T) {
	tests := []struct {
		model string
		want  bool
	}{
		{"gpt-4", true},
		{"gpt-3.5-turbo", true},
		{"claude-3-opus", true},
		{"claude-4.1-opus", true},
		{"gemini-1.5-flash", false},
		{"", true}, // Empty defaults to tiktoken
	}
	
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			if got := IsTiktokenModel(tt.model); got != tt.want {
				t.Errorf("IsTiktokenModel(%q) = %v, want %v", tt.model, got, tt.want)
			}
		})
	}
}

func TestDefaultTokenizerFactory(t *testing.T) {
	factory := &DefaultTokenizerFactory{}
	
	tests := []struct {
		model     string
		wantError bool
	}{
		{"gpt-4", false},
		{"claude-3-opus", false},
		// Note: Gemini will fail without proper setup, but that's expected
		{"", true}, // Empty model should error
	}
	
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			_, err := factory.CreateTokenizer(tt.model)
			if (err != nil) != tt.wantError {
				t.Errorf("CreateTokenizer(%q) error = %v, wantError %v", tt.model, err, tt.wantError)
			}
		})
	}
}

func TestGetSupportedModels(t *testing.T) {
	models := GetSupportedModels()
	if len(models) == 0 {
		t.Error("GetSupportedModels() should return non-empty list")
	}
	
	// Check that both Gemini and Tiktoken models are included
	hasGemini := false
	hasTiktoken := false
	
	for _, model := range models {
		if strings.Contains(model, "gemini") {
			hasGemini = true
		}
		if strings.Contains(model, "gpt") || strings.Contains(model, "claude") {
			hasTiktoken = true
		}
	}
	
	if !hasGemini {
		t.Error("GetSupportedModels() should include Gemini models")
	}
	if !hasTiktoken {
		t.Error("GetSupportedModels() should include OpenAI/Claude models")
	}
}

func TestIsModelSupported(t *testing.T) {
	tests := []struct {
		model string
		want  bool
	}{
		{"gpt-4", true},
		{"GPT-4", true}, // Case insensitive
		{"claude-3-opus", true},
		{"gemini-1.5-flash", true},
		{"unknown-model", false},
		{"", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.model, func(t *testing.T) {
			if got := IsModelSupported(tt.model); got != tt.want {
				t.Errorf("IsModelSupported(%q) = %v, want %v", tt.model, got, tt.want)
			}
		})
	}
}