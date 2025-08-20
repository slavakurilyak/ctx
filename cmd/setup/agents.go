package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewAgentsCmd creates the setup agents subcommand for AGENTS.md standard
func NewAgentsCmd() *cobra.Command {
	var force bool

	agentsCmd := &cobra.Command{
		Use:   "agents",
		Short: "Generate AGENTS.md with ctx documentation (OpenAI standard)",
		Long: `Generate an AGENTS.md file with ctx documentation following the OpenAI AGENTS.md standard.

AGENTS.md is an open format for guiding coding agents, supported by:
- OpenAI Codex
- Cursor
- RooCode  
- OpenCode
- Amp
- Jules (Google)
- Factory
- And many others

This creates a standardized AGENTS.md file that works across multiple AI coding assistants.

Examples:
  ctx setup agents         # Generate AGENTS.md
  ctx setup agents --force # Overwrite existing AGENTS.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := "AGENTS.md"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			content, err := GenerateRules("AGENTS.md")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			fmt.Println("\nThis AGENTS.md file works with OpenAI Codex, Cursor, RooCode, and other compatible agents.")
			fmt.Println("Learn more about the AGENTS.md standard at https://agents.md/")
			return nil
		},
	}

	agentsCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return agentsCmd
}