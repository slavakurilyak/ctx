package tokenizer

import (
	"fmt"
	"sync"
)

// TokenizerCache manages a cache of tokenizers to avoid repeated initialization
type TokenizerCache struct {
	factory TokenizerFactory
	cache   map[string]Tokenizer
	mu      sync.RWMutex
}

// NewTokenizerCache creates a new tokenizer cache with the given factory
func NewTokenizerCache(factory TokenizerFactory) *TokenizerCache {
	if factory == nil {
		factory = &DefaultTokenizerFactory{}
	}

	return &TokenizerCache{
		factory: factory,
		cache:   make(map[string]Tokenizer),
	}
}

// GetOrCreate returns a cached tokenizer or creates a new one if not found
func (c *TokenizerCache) GetOrCreate(modelName string) (Tokenizer, error) {
	if modelName == "" {
		return nil, fmt.Errorf("model name cannot be empty")
	}

	// Try to get from cache with read lock
	c.mu.RLock()
	if tok, exists := c.cache[modelName]; exists {
		c.mu.RUnlock()
		return tok, nil
	}
	c.mu.RUnlock()

	// Not in cache, acquire write lock to create
	c.mu.Lock()
	defer c.mu.Unlock()

	// Double-check after acquiring write lock
	if tok, exists := c.cache[modelName]; exists {
		return tok, nil
	}

	// Create new tokenizer
	tok, err := c.factory.CreateTokenizer(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to create tokenizer for model %s: %w", modelName, err)
	}

	// Cache the tokenizer
	c.cache[modelName] = tok
	return tok, nil
}

// Clear removes all cached tokenizers
func (c *TokenizerCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Clear the map
	for k := range c.cache {
		delete(c.cache, k)
	}
}

// Size returns the number of cached tokenizers
func (c *TokenizerCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}

// Has checks if a tokenizer for the given model is cached
func (c *TokenizerCache) Has(modelName string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.cache[modelName]
	return exists
}

// DefaultCache is a global cache instance for convenience
var DefaultCache = NewTokenizerCache(nil)

// GetCachedTokenizer is a convenience function using the default cache
func GetCachedTokenizer(modelName string) (Tokenizer, error) {
	return DefaultCache.GetOrCreate(modelName)
}
