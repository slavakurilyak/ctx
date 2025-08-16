package tokenizer

// Tokenizer is the common interface for all tokenizer implementations
type Tokenizer interface {
	// CountTokens counts the number of tokens in the given text
	CountTokens(text string) (int, error)
	
	// GetModelName returns the name of the model this tokenizer is configured for
	GetModelName() string
}

// TokenizerFactory creates tokenizers based on model names
type TokenizerFactory interface {
	// CreateTokenizer creates a tokenizer for the specified model
	CreateTokenizer(modelName string) (Tokenizer, error)
}