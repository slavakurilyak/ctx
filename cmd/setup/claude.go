package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewClaudeCmd creates the setup claude subcommand
func NewClaudeCmd() *cobra.Command {
	var force bool

	claudeCmd := &cobra.Command{
		Use:   "claude",
		Short: "Generate CLAUDE.md with ctx documentation",
		Long: `Generate a CLAUDE.md file with ctx documentation for Claude AI assistant.

This command creates a CLAUDE.md file in the current directory
with comprehensive ctx documentation formatted for Claude.

Examples:
  ctx setup claude         # Generate CLAUDE.md
  ctx setup claude --force # Overwrite existing CLAUDE.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := "CLAUDE.md"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			markdownContent, err := GenerateRules("CLAUDE.md")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			return nil
		},
	}

	claudeCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return claudeCmd
}
