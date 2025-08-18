package telemetry

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/slavakurilyak/ctx/internal/version"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

const (
	serviceName = "ctx"
)

type Manager struct {
	enabled        bool
	tracer         trace.Tracer
	tracerProvider *sdktrace.TracerProvider
}

var (
	globalManager *Manager
)

// Initialize sets up OpenTelemetry with OTLP exporter
func Initialize(ctx context.Context) (*Manager, error) {
	// Check if telemetry is disabled
	if isDisabled() {
		return &Manager{enabled: false}, nil
	}

	// Create resource describing this application
	res, err := createResource(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP exporter
	exporter, err := createOTLPExporter(ctx)
	if err != nil {
		// If exporter fails, telemetry is disabled but app continues
		return &Manager{enabled: false}, nil
	}

	// Create TracerProvider with the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Register as global tracer provider
	otel.SetTracerProvider(tp)

	// Set up propagation (W3C Trace Context by default)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	manager := &Manager{
		enabled:        true,
		tracer:         tp.Tracer("github.com/slavakurilyak/ctx"),
		tracerProvider: tp,
	}

	globalManager = manager
	return manager, nil
}

// GetGlobal returns the global telemetry manager
func GetGlobal() *Manager {
	if globalManager == nil {
		return &Manager{enabled: false}
	}
	return globalManager
}

// StartCommandSpan starts a new span for command execution
func (m *Manager) StartCommandSpan(ctx context.Context, command string) (context.Context, trace.Span) {
	if !m.enabled || m.tracer == nil {
		// Return a no-op span
		return ctx, trace.SpanFromContext(ctx)
	}

	return m.tracer.Start(ctx, "command.execute",
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(
			attribute.String("code.function", "execute"),
			attribute.String("code.namespace", "ctx"),
			attribute.String("command.text", command),
			attribute.String("command.type", detectCommandType(command)),
		),
	)
}

// RecordCommandResult records the result of a command execution
func (m *Manager) RecordCommandResult(span trace.Span, exitCode int, duration time.Duration, outputBytes int, tokenCount *int) {
	if !m.enabled || span == nil {
		return
	}

	// Add result attributes
	span.SetAttributes(
		attribute.Int("command.exit_code", exitCode),
		attribute.Int64("command.duration_ms", duration.Milliseconds()),
		attribute.Int("command.output_bytes", outputBytes),
	)

	if tokenCount != nil {
		span.SetAttributes(attribute.Int("command.token_count", *tokenCount))
	}

	// Set status based on exit code
	if exitCode != 0 {
		span.SetAttributes(attribute.Bool("command.error", true))
	}
}

// AddSpanAttributes adds custom attributes to a span
func (m *Manager) AddSpanAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	if !m.enabled || span == nil {
		return
	}
	span.SetAttributes(attrs...)
}

// GetTraceContext extracts trace context as a map for storage
func (m *Manager) GetTraceContext(ctx context.Context) map[string]string {
	if !m.enabled {
		return nil
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return nil
	}

	spanCtx := span.SpanContext()
	if !spanCtx.IsValid() {
		return nil
	}

	return map[string]string{
		"trace_id":    spanCtx.TraceID().String(),
		"span_id":     spanCtx.SpanID().String(),
		"trace_flags": fmt.Sprintf("%02x", spanCtx.TraceFlags()),
	}
}

// Shutdown gracefully shuts down the telemetry system
func (m *Manager) Shutdown(ctx context.Context) error {
	if !m.enabled || m.tracerProvider == nil {
		return nil
	}

	return m.tracerProvider.Shutdown(ctx)
}

// IsEnabled returns whether telemetry is enabled
func (m *Manager) IsEnabled() bool {
	return m.enabled
}

// Helper functions

func isDisabled() bool {
	// Check various environment variables that might disable telemetry
	disableVars := []string{
		"CTX_NO_TELEMETRY",
		"CTX_NO_OPENTELEMETRY",
		"CTX_PRIVATE",
		"OTEL_SDK_DISABLED",
	}

	for _, v := range disableVars {
		if os.Getenv(v) == "true" || os.Getenv(v) == "1" {
			return true
		}
	}

	return false
}

func createResource(ctx context.Context) (*resource.Resource, error) {
	// Get hostname
	hostname, _ := os.Hostname()

	// Get user
	user := os.Getenv("USER")
	if user == "" {
		user = os.Getenv("USERNAME")
	}

	// Create resource with attributes
	return resource.NewWithAttributes(
		"",
		attribute.String("service.name", serviceName),
		attribute.String("service.version", version.GetVersion()),
		attribute.String("host.name", hostname),
		attribute.String("user.name", user),
		attribute.String("os.type", getOSType()),
		attribute.String("process.runtime.name", "go"),
		attribute.String("process.runtime.version", getGoVersion()),
	), nil
}

func createOTLPExporter(ctx context.Context) (*otlptrace.Exporter, error) {
	// Get OTLP endpoint from environment or use default
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = os.Getenv("CTX_OTLP_ENDPOINT")
	}
	if endpoint == "" {
		endpoint = "localhost:4318" // Default OTLP HTTP endpoint
	}

	// Configure OTLP exporter options
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}

	// Check for insecure mode (for local development)
	if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true" ||
		os.Getenv("CTX_OTLP_INSECURE") == "true" {
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	// Add headers if specified
	if headers := os.Getenv("OTEL_EXPORTER_OTLP_HEADERS"); headers != "" {
		opts = append(opts, otlptracehttp.WithHeaders(parseHeaders(headers)))
	}

	// Create the exporter
	client := otlptracehttp.NewClient(opts...)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	return exporter, nil
}

func detectCommandType(command string) string {
	// Simple command type detection based on the first word
	if command == "" {
		return "unknown"
	}

	// Extract first word
	firstSpace := -1
	for i, ch := range command {
		if ch == ' ' {
			firstSpace = i
			break
		}
	}

	cmdName := command
	if firstSpace > 0 {
		cmdName = command[:firstSpace]
	}

	// Categorize by command type
	switch cmdName {
	case "psql", "mysql", "sqlite3", "redis-cli", "mongosh":
		return "database"
	case "docker", "kubectl", "docker-compose":
		return "container"
	case "git", "svn", "hg":
		return "vcs"
	case "npm", "yarn", "pnpm", "pip", "go", "cargo", "mvn", "gradle":
		return "build"
	case "curl", "wget", "http", "httpie":
		return "http"
	case "ls", "cat", "grep", "find", "sed", "awk", "sort", "uniq":
		return "filesystem"
	case "ps", "top", "htop", "df", "du", "free", "uptime":
		return "system"
	default:
		return "general"
	}
}

func parseHeaders(headers string) map[string]string {
	result := make(map[string]string)
	// Simple header parsing (key=value,key2=value2)
	pairs := splitByComma(headers)
	for _, pair := range pairs {
		if idx := indexOfEqual(pair); idx > 0 {
			key := pair[:idx]
			value := pair[idx+1:]
			result[key] = value
		}
	}
	return result
}

func splitByComma(s string) []string {
	var result []string
	var current string
	for _, ch := range s {
		if ch == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func indexOfEqual(s string) int {
	for i, ch := range s {
		if ch == '=' {
			return i
		}
	}
	return -1
}

func getOSType() string {
	// Simple OS detection
	if os.Getenv("OS") == "Windows_NT" {
		return "windows"
	}
	// Could use runtime.GOOS for more accuracy
	return "unix"
}

func getGoVersion() string {
	// Would normally use runtime.Version()
	return "1.21+"
}
