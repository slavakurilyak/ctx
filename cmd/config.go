package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/slavakurilyak/ctx/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewConfigCmd creates the config command with subcommands
func NewConfigCmd() *cobra.Command {
	configCmd := &cobra.Command{
		Use:   "config",
		Short: "Manage ctx configuration",
		Long:  "View and manage ctx configuration settings",
	}

	configCmd.AddCommand(newConfigViewCmd())

	return configCmd
}

// newConfigViewCmd creates the config view subcommand
func newConfigViewCmd() *cobra.Command {
	var outputFormat string

	viewCmd := &cobra.Command{
		Use:   "view",
		Short: "View the current configuration",
		Long: `View the current configuration with all settings merged from:
1. Default values
2. Configuration file (~/.config/ctx/config.yaml)
3. Environment variables
4. Command-line flags

The output shows the final configuration and the source of each setting.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get the current configuration
			cfg := config.NewFromFlagsAndEnv(cmd)

			// Build the configuration output
			configOutput := buildConfigOutput(cfg)

			// Output in requested format
			switch outputFormat {
			case "json":
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				return encoder.Encode(configOutput)
			case "yaml":
				encoder := yaml.NewEncoder(os.Stdout)
				encoder.SetIndent(2)
				return encoder.Encode(configOutput)
			default:
				// Human-readable format
				printHumanReadable(configOutput)
				return nil
			}
		},
	}

	viewCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Output format: text, json, or yaml")

	return viewCmd
}

// ConfigOutput represents the configuration with source information
type ConfigOutput struct {
	TokenModel        ConfigValue   `json:"token_model" yaml:"token_model"`
	DefaultTimeout    ConfigValue   `json:"default_timeout" yaml:"default_timeout"`
	OutputFormat      ConfigValue   `json:"output_format" yaml:"output_format"`
	CacheDir          ConfigValue   `json:"cache_dir" yaml:"cache_dir"`
	NoTokens          ConfigValue   `json:"no_tokens" yaml:"no_tokens"`
	NoHistory         ConfigValue   `json:"no_history" yaml:"no_history"`
	NoTelemetry       ConfigValue   `json:"no_telemetry" yaml:"no_telemetry"`
	Limits            LimitsOutput  `json:"limits" yaml:"limits"`
}

// LimitsOutput represents the limits configuration with source information
type LimitsOutput struct {
	MaxTokens         ConfigValue `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
	MaxOutputBytes    ConfigValue `json:"max_output_bytes,omitempty" yaml:"max_output_bytes,omitempty"`
	MaxLines          ConfigValue `json:"max_lines,omitempty" yaml:"max_lines,omitempty"`
	MaxPipelineStages ConfigValue `json:"max_pipeline_stages,omitempty" yaml:"max_pipeline_stages,omitempty"`
}

// ConfigValue represents a configuration value with its source
type ConfigValue struct {
	Value  interface{} `json:"value" yaml:"value"`
	Source string      `json:"source" yaml:"source"`
}

// buildConfigOutput builds the configuration output with source information
func buildConfigOutput(cfg *config.Config) ConfigOutput {
	output := ConfigOutput{
		TokenModel: ConfigValue{
			Value:  cfg.TokenModel,
			Source: getSource(cfg, "TokenModel"),
		},
		DefaultTimeout: ConfigValue{
			Value:  cfg.DefaultTimeout.String(),
			Source: getSource(cfg, "DefaultTimeout"),
		},
		OutputFormat: ConfigValue{
			Value:  cfg.OutputFormat,
			Source: getSource(cfg, "OutputFormat"),
		},
		CacheDir: ConfigValue{
			Value:  cfg.CacheDir,
			Source: getSource(cfg, "CacheDir"),
		},
		NoTokens: ConfigValue{
			Value:  cfg.NoTokens,
			Source: getSource(cfg, "NoTokens"),
		},
		NoHistory: ConfigValue{
			Value:  cfg.NoHistory,
			Source: getSource(cfg, "NoHistory"),
		},
		NoTelemetry: ConfigValue{
			Value:  cfg.NoTelemetry,
			Source: cfg.NoTelemetrySource,
		},
		Limits: LimitsOutput{},
	}

	// Handle limits
	if cfg.Limits.MaxTokens != nil {
		output.Limits.MaxTokens = ConfigValue{
			Value:  *cfg.Limits.MaxTokens,
			Source: getSource(cfg, "Limits.MaxTokens"),
		}
	}
	if cfg.Limits.MaxOutputBytes != nil {
		output.Limits.MaxOutputBytes = ConfigValue{
			Value:  *cfg.Limits.MaxOutputBytes,
			Source: getSource(cfg, "Limits.MaxOutputBytes"),
		}
	}
	if cfg.Limits.MaxLines != nil {
		output.Limits.MaxLines = ConfigValue{
			Value:  *cfg.Limits.MaxLines,
			Source: getSource(cfg, "Limits.MaxLines"),
		}
	}
	if cfg.Limits.MaxPipelineStages != nil {
		output.Limits.MaxPipelineStages = ConfigValue{
			Value:  *cfg.Limits.MaxPipelineStages,
			Source: getSource(cfg, "Limits.MaxPipelineStages"),
		}
	}

	return output
}

// getSource determines the source of a configuration value
// This is a simplified implementation - in a real system, you'd track
// the actual source during configuration loading
func getSource(cfg *config.Config, field string) string {
	// Check if environment variable is set
	envVars := map[string]string{
		"TokenModel":            "CTX_TOKEN_MODEL",
		"DefaultTimeout":        "CTX_TIMEOUT",
		"OutputFormat":          "CTX_OUTPUT_FORMAT",
		"NoTokens":              "CTX_NO_TOKENS",
		"NoHistory":             "CTX_NO_HISTORY",
		"Limits.MaxTokens":      "CTX_MAX_TOKENS",
		"Limits.MaxOutputBytes": "CTX_MAX_OUTPUT_BYTES",
		"Limits.MaxLines":       "CTX_MAX_LINES",
		"Limits.MaxPipelineStages": "CTX_MAX_PIPELINE_STAGES",
	}

	if envVar, ok := envVars[field]; ok {
		if os.Getenv(envVar) != "" {
			return "environment variable"
		}
	}

	// Check if value is non-zero (likely set from file or flag)
	v := reflect.ValueOf(cfg).Elem()
	fieldValue := v.FieldByName(field)
	
	if field == "NoTelemetry" {
		return cfg.NoTelemetrySource
	}

	if fieldValue.IsValid() && !fieldValue.IsZero() {
		// Could be from file or flag - would need more tracking to differentiate
		return "config file or flag"
	}

	return "default"
}

// printHumanReadable prints the configuration in a human-readable format
func printHumanReadable(output ConfigOutput) {
	fmt.Println("Current ctx Configuration:")
	fmt.Println("==========================")
	fmt.Printf("Token Model:        %v (source: %s)\n", output.TokenModel.Value, output.TokenModel.Source)
	fmt.Printf("Default Timeout:    %v (source: %s)\n", output.DefaultTimeout.Value, output.DefaultTimeout.Source)
	fmt.Printf("Output Format:      %v (source: %s)\n", output.OutputFormat.Value, output.OutputFormat.Source)
	fmt.Printf("Cache Dir:          %v (source: %s)\n", output.CacheDir.Value, output.CacheDir.Source)
	fmt.Printf("No Tokens:          %v (source: %s)\n", output.NoTokens.Value, output.NoTokens.Source)
	fmt.Printf("No History:         %v (source: %s)\n", output.NoHistory.Value, output.NoHistory.Source)
	fmt.Printf("No Telemetry:       %v (source: %s)\n", output.NoTelemetry.Value, output.NoTelemetry.Source)
	
	fmt.Println("\nLimits:")
	fmt.Println("-------")
	if output.Limits.MaxTokens.Value != nil {
		fmt.Printf("  Max Tokens:       %v (source: %s)\n", output.Limits.MaxTokens.Value, output.Limits.MaxTokens.Source)
	} else {
		fmt.Println("  Max Tokens:       not set")
	}
	if output.Limits.MaxOutputBytes.Value != nil {
		fmt.Printf("  Max Output Bytes: %v (source: %s)\n", output.Limits.MaxOutputBytes.Value, output.Limits.MaxOutputBytes.Source)
	} else {
		fmt.Println("  Max Output Bytes: not set")
	}
	if output.Limits.MaxLines.Value != nil {
		fmt.Printf("  Max Lines:        %v (source: %s)\n", output.Limits.MaxLines.Value, output.Limits.MaxLines.Source)
	} else {
		fmt.Println("  Max Lines:        not set")
	}
	if output.Limits.MaxPipelineStages.Value != nil {
		fmt.Printf("  Max Pipeline Stages: %v (source: %s)\n", output.Limits.MaxPipelineStages.Value, output.Limits.MaxPipelineStages.Source)
	} else {
		fmt.Println("  Max Pipeline Stages: not set")
	}
}