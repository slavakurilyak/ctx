package models

import (
	"time"
)

// CurrentSchemaVersion defines the version of the JSON output schema.
// Pre-1.0 versions (0.x) indicate the schema is still evolving.
// This should be updated when making changes to the output structure:
// - Increment MINOR (e.g., 0.1 -> 0.2) for any changes before 1.0
// - Version 1.0 will indicate a stable, backward-compatible schema
// See docs/VERSIONING.md for detailed versioning strategy.
const CurrentSchemaVersion = "0.1"

// Output represents the final output structure for ctx commands
// This is the structure that gets printed to console and saved to history
type Output struct {
	Tokens        int              `json:"tokens"`         // Token count - most important, shown first
	Output        string           `json:"output"`         // Command output - second most important
	Input         string           `json:"input"`          // Command executed - third
	Metadata      MetadataSection  `json:"metadata"`       // Additional details
	Telemetry     *TelemetrySection `json:"telemetry,omitempty"`
	SchemaVersion string           `json:"schema_version"` // Schema version for parsers
}

// MetadataSection contains execution details and context
type MetadataSection struct {
	// Execution results
	ExitCode      int    `json:"exit_code"`
	Success       bool   `json:"success"` // Derived from exit_code (0 = true)
	Error         string `json:"error,omitempty"` // Error message if one occurred
	FailureReason string `json:"failure_reason,omitempty"` // Machine-readable failure reason (e.g., "line_limit_exceeded")
	
	// Performance metrics
	Duration int `json:"duration"` // Execution time in milliseconds
	Bytes    int `json:"bytes"`    // Output size in bytes
	
	// Context information
	Timestamp string `json:"timestamp"`  // RFC3339 formatted timestamp
	Directory string `json:"directory"`  // Working directory
	User      string `json:"user"`       // Username
	Host      string `json:"host"`       // Hostname
	SessionID string `json:"session_id"` // Unique session identifier
}

// TelemetrySection contains optional OpenTelemetry trace information
type TelemetrySection struct {
	TraceID    string `json:"trace_id"`    // Distributed trace identifier
	SpanID     string `json:"span_id"`     // Span identifier
	TraceFlags string `json:"trace_flags"` // W3C trace flags
}

// StreamEvent represents a streaming event for long-running commands
type StreamEvent struct {
	Type     string  `json:"type"`               // "stdout", "stderr", or "result"
	Line     string  `json:"line,omitempty"`     // The line of output for stdout/stderr events
	Envelope *Output `json:"envelope,omitempty"` // The final envelope for the result event
}

// NewOutput creates a new Output structure
func NewOutput(command string, output []byte, exitCode int, duration time.Duration) *Output {
	return &Output{
		Tokens:        0, // Will be set by enricher
		Output:        string(output),
		Input:         command,
		SchemaVersion: CurrentSchemaVersion,
		Metadata: MetadataSection{
			ExitCode: exitCode,
			Success:  exitCode == 0,
			Duration: int(duration.Milliseconds()),
			Bytes:    len(output),
			// Other fields will be populated by enricher
		},
	}
}