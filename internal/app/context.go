package app

import (
	"github.com/slavakurilyak/ctx/internal/config"
	"github.com/slavakurilyak/ctx/internal/history"
	"github.com/slavakurilyak/ctx/internal/telemetry"
	"github.com/slavakurilyak/ctx/internal/tokenizer"
)

// Add a context key for type-safe context value access.
type contextKey string

const AppContextKey contextKey = "appContext"

// AppContext is the central dependency injection container for the application
type AppContext struct {
	// Tokenizer handles token counting for different LLM models
	Tokenizer tokenizer.Tokenizer
	
	// TokenizerCache manages tokenizer instances
	TokenizerCache *tokenizer.TokenizerCache
	
	// Telemetry manages OpenTelemetry tracing
	Telemetry *telemetry.Manager
	
	// History manages command history
	History *history.HistoryManager
	
	// Config holds application configuration
	Config *config.Config
}

// Option is a functional option for configuring AppContext
type Option func(*AppContext)

// NewAppContext creates a new application context with the given options
func NewAppContext(opts ...Option) *AppContext {
	ctx := &AppContext{
		Config: config.Get(),
	}
	
	// Apply all options
	for _, opt := range opts {
		opt(ctx)
	}
	
	// Set defaults if not provided
	if ctx.History == nil {
		ctx.History = history.NewHistoryManager()
	}
	
	if ctx.TokenizerCache == nil {
		ctx.TokenizerCache = tokenizer.NewTokenizerCache(nil)
	}
	
	return ctx
}

// WithTokenizer sets a specific tokenizer
func WithTokenizer(tok tokenizer.Tokenizer) Option {
	return func(ctx *AppContext) {
		ctx.Tokenizer = tok
	}
}

// WithTokenizerCache sets the tokenizer cache
func WithTokenizerCache(cache *tokenizer.TokenizerCache) Option {
	return func(ctx *AppContext) {
		ctx.TokenizerCache = cache
	}
}

// WithTelemetry sets the telemetry manager
func WithTelemetry(tm *telemetry.Manager) Option {
	return func(ctx *AppContext) {
		ctx.Telemetry = tm
	}
}

// WithHistory sets the history manager
func WithHistory(hm *history.HistoryManager) Option {
	return func(ctx *AppContext) {
		ctx.History = hm
	}
}

// WithConfig sets the configuration
func WithConfig(cfg *config.Config) Option {
	return func(ctx *AppContext) {
		ctx.Config = cfg
	}
}

// GetTokenizer returns the configured tokenizer, creating one if necessary
func (ctx *AppContext) GetTokenizer() (tokenizer.Tokenizer, error) {
	if ctx.Tokenizer != nil {
		return ctx.Tokenizer, nil
	}
	
	// Get model from config
	model := ctx.Config.TokenModel
	if model == "" {
		model = "claude-4.1-opus" // Default
	}
	
	// Get or create from cache
	tok, err := ctx.TokenizerCache.GetOrCreate(model)
	if err != nil {
		return nil, err
	}
	
	ctx.Tokenizer = tok
	return tok, nil
}