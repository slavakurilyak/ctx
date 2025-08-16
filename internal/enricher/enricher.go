package enricher

import (
	"context"
	"os"
	"strings"
	"time"
	
	"github.com/google/uuid"
	"github.com/slavakurilyak/ctx/internal/config"
	"github.com/slavakurilyak/ctx/internal/executor"
	"github.com/slavakurilyak/ctx/internal/history"
	"github.com/slavakurilyak/ctx/internal/models"
	"github.com/slavakurilyak/ctx/internal/telemetry"
	"github.com/slavakurilyak/ctx/internal/tokenizer"
)

// Enricher handles output enrichment with injected dependencies
type Enricher struct {
	tokenizer tokenizer.Tokenizer
	history   *history.HistoryManager
	telemetry *telemetry.Manager
	config    *config.Config
}

// NewEnricher creates a new enricher with the given dependencies
func NewEnricher(tok tokenizer.Tokenizer, hist *history.HistoryManager, tel *telemetry.Manager, cfg *config.Config) *Enricher {
	return &Enricher{
		tokenizer: tok,
		history:   hist,
		telemetry: tel,
		config:    cfg,
	}
}

// EnrichOutput enriches the execution result with metadata and token counts
func (e *Enricher) EnrichOutput(ctx context.Context, result *executor.ExecutionResult) (*models.Output, error) {
	output := models.NewOutput(
		result.Command,
		result.Output,
		result.ExitCode,
		result.Duration,
	)
	
	// Populate metadata context fields
	output.Metadata.Timestamp = time.Now().Format(time.RFC3339)
	output.Metadata.SessionID = uuid.New().String()
	
	// Get working directory
	if cwd, err := os.Getwd(); err == nil {
		output.Metadata.Directory = cwd
	}
	
	// Get user
	if user := os.Getenv("USER"); user != "" {
		output.Metadata.User = user
	}
	
	// Get hostname
	if hostname, err := os.Hostname(); err == nil {
		output.Metadata.Host = hostname
	}
	
	// Count tokens if enabled and tokenizer is available
	if e.shouldCountTokens() && e.tokenizer != nil {
		tokenCount, err := e.tokenizer.CountTokens(string(result.Output))
		if err == nil {
			output.Tokens = tokenCount
		}
		// If token counting fails, we still return the output without tokens
	}
	
	// Get trace context if telemetry is enabled
	if e.telemetry != nil {
		if traceContext := e.telemetry.GetTraceContext(ctx); traceContext != nil {
			output.Telemetry = &models.TelemetrySection{
				TraceID:    traceContext["trace_id"],
				SpanID:     traceContext["span_id"],
				TraceFlags: traceContext["trace_flags"],
			}
		}
	}
	
	// Save to history if history manager is available
	if e.history != nil {
		e.history.SaveRecord(output)
	}
	
	return output, nil
}

// shouldCountTokens checks if token counting is enabled
func (e *Enricher) shouldCountTokens() bool {
	// Check if token counting is disabled in config
	if e.config != nil && e.config.NoTokens {
		return false
	}
	return true
}

// ShouldEnrich checks if a command should be enriched
func ShouldEnrich(command string) bool {
	// Don't double-wrap ctx commands
	cmdParts := strings.Fields(command)
	if len(cmdParts) > 0 && strings.HasSuffix(cmdParts[0], "ctx") {
		return false
	}
	return true
}

// Option is a functional option for configuring the Enricher
type Option func(*Enricher)

// WithTokenizer sets the tokenizer
func WithTokenizer(tok tokenizer.Tokenizer) Option {
	return func(e *Enricher) {
		e.tokenizer = tok
	}
}

// WithHistory sets the history manager
func WithHistory(hist *history.HistoryManager) Option {
	return func(e *Enricher) {
		e.history = hist
	}
}

// WithTelemetry sets the telemetry manager
func WithTelemetry(tel *telemetry.Manager) Option {
	return func(e *Enricher) {
		e.telemetry = tel
	}
}

// NewEnricherWithOptions creates a new enricher with functional options
func NewEnricherWithOptions(opts ...Option) *Enricher {
	e := &Enricher{}
	for _, opt := range opts {
		opt(e)
	}
	return e
}