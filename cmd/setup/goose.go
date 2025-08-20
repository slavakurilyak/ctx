package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewGooseCmd creates the setup goose subcommand
func NewGooseCmd() *cobra.Command {
	var force bool

	gooseCmd := &cobra.Command{
		Use:   "goose",
		Short: "Generate GOOSE.md with ctx documentation",
		Long: `Generate a GOOSE.md file with ctx documentation for Goose AI agent.

This command creates a GOOSE.md file in the current directory
with comprehensive ctx documentation formatted for Goose.

Goose uses GOOSE.md files for system prompts and context,
supporting MCP servers and extensible configurations.

Examples:
  ctx setup goose         # Generate GOOSE.md
  ctx setup goose --force # Overwrite existing GOOSE.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := "GOOSE.md"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			markdownContent, err := GenerateRules("GOOSE.md")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			fmt.Println("\nGoose will use this context for all sessions in this project.")
			fmt.Println("For global context, place this file at ~/.goose/GOOSE.md")
			return nil
		},
	}

	gooseCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return gooseCmd
}
