package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewGeminiCmd creates the setup gemini subcommand
func NewGeminiCmd() *cobra.Command {
	var force bool

	geminiCmd := &cobra.Command{
		Use:   "gemini",
		Short: "Generate GEMINI.md with ctx documentation",
		Long: `Generate a GEMINI.md file with ctx documentation for Gemini CLI.

This command creates a GEMINI.md file in the current directory
with comprehensive ctx documentation formatted for Gemini CLI.

Gemini CLI uses GEMINI.md files for system prompts and context,
supporting MCP (Model Context Protocol) and extensible configurations.

Examples:
  ctx setup gemini         # Generate GEMINI.md
  ctx setup gemini --force # Overwrite existing GEMINI.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := "GEMINI.md"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Get ctx help output
			content, err := GetCtxHelp()
			if err != nil {
				return err
			}

			// Format as markdown
			markdownContent := fmt.Sprintf(`# ctx Documentation

%s
`, content)

			// Write the file
			if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			fmt.Println("\nGemini CLI will use this context for all prompts in this project.")
			fmt.Println("For global context, place this file at ~/.gemini/GEMINI.md")
			return nil
		},
	}

	geminiCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return geminiCmd
}
