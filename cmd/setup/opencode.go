package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewOpenCodeCmd creates the setup opencode subcommand
func NewOpenCodeCmd() *cobra.Command {
	var force bool

	openCodeCmd := &cobra.Command{
		Use:   "opencode",
		Short: "Generate AGENTS.md with ctx documentation",
		Long: `Generate an AGENTS.md file with ctx documentation for OpenCode AI assistant.

This command creates an AGENTS.md file in the current directory
with ctx help output for OpenCode to understand ctx usage.

Examples:
  ctx setup opencode         # Generate AGENTS.md
  ctx setup opencode --force # Overwrite existing AGENTS.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := "AGENTS.md"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			// OpenCode uses AGENTS.md format
			content, err := GenerateRules("AGENTS.md")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			return nil
		},
	}

	openCodeCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return openCodeCmd
}