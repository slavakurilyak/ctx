package config

import (
	"fmt"
	"strings"
	"github.com/charmbracelet/lipgloss"
)

// EnvVar represents a single environment variable with metadata
type EnvVar struct {
	Name        string
	Description string
	Example     string // Optional example value
}

// EnvVars defines all environment variables with their documentation
var EnvVars = []EnvVar{
	{
		Name:        "CTX_TOKEN_MODEL",
		Description: "Sets the default token counting provider",
		Example:     "\"anthropic\", \"openai\", \"gemini\"",
	},
	{
		Name:        "CTX_NO_TOKENS",
		Description: "If \"true\", disables token counting for all commands",
	},
	{
		Name:        "CTX_PRETTY",
		Description: "If \"true\", outputs in pretty format instead of JSON",
	},
	{
		Name:        "CTX_MAX_TOKENS",
		Description: "Sets the maximum number of tokens allowed in output",
		Example:     "\"5000\"",
	},
	{
		Name:        "CTX_MAX_OUTPUT_BYTES",
		Description: "Sets the maximum number of bytes allowed in output",
		Example:     "\"1048576\" for 1MB",
	},
	{
		Name:        "CTX_MAX_LINES",
		Description: "Sets the maximum number of lines allowed in output",
		Example:     "\"1000\"",
	},
	{
		Name:        "CTX_MAX_PIPELINE_STAGES",
		Description: "Sets the maximum number of pipeline stages allowed",
		Example:     "\"5\"",
	},
	{
		Name:        "CTX_NO_HISTORY",
		Description: "If \"true\", disables command history recording",
	},
	{
		Name:        "CTX_NO_TELEMETRY",
		Description: "If \"true\", disables OpenTelemetry tracing",
	},
	{
		Name:        "CTX_PRIVATE",
		Description: "If \"true\", is equivalent to setting both CTX_NO_HISTORY and CTX_NO_TELEMETRY to \"true\"",
	},
	{
		Name:        "CTX_TIMEOUT",
		Description: "Sets a default timeout for commands",
		Example:     "\"30s\", \"1m\"",
	},
	{
		Name:        "CTX_API_ENDPOINT",
		Description: "API endpoint for ctx Pro features (required for Pro features)",
		Example:     "\"https://api.ctx.click\"",
	},
	{
		Name:        "OTEL_EXPORTER_OTLP_ENDPOINT",
		Description: "The OTLP endpoint for telemetry data",
		Example:     "\"http://localhost:4318\"",
	},
}

// Styling for environment section
var (
	envHeaderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B35")).
		Bold(true)
	
	envVarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFA500")).
		Bold(true)
	
	envDescStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#555555"))
)

// GenerateEnvSection creates the ENVIRONMENT section for help text
func GenerateEnvSection() string {
	var sb strings.Builder
	sb.WriteString(envHeaderStyle.Render("ENVIRONMENT:") + "\n\n")
	sb.WriteString("  " + envDescStyle.Render("CLI flags take precedence over environment variables.") + "\n\n")
	
	for _, env := range EnvVars {
		sb.WriteString(fmt.Sprintf("  %s\n", envVarStyle.Render(env.Name)))
		sb.WriteString(fmt.Sprintf("    %s", envDescStyle.Render(env.Description)))
		if env.Example != "" {
			sb.WriteString(fmt.Sprintf(" %s", envDescStyle.Render(fmt.Sprintf("(e.g., %s)", env.Example))))
		}
		sb.WriteString(".\n\n")
	}
	
	// Remove the extra newline at the end
	result := sb.String()
	return strings.TrimSuffix(result, "\n")
}