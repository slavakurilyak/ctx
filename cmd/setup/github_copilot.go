package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewGitHubCopilotCmd creates the setup github-copilot subcommand
func NewGitHubCopilotCmd() *cobra.Command {
	var force bool

	copilotCmd := &cobra.Command{
		Use:   "github-copilot",
		Short: "Generate GitHub Copilot instructions with ctx documentation",
		Long: `Generate GitHub Copilot instructions file with ctx documentation.

This command creates .github/copilot-instructions.md with comprehensive
ctx documentation formatted for GitHub Copilot.

Examples:
  ctx setup github-copilot         # Generate .github/copilot-instructions.md
  ctx setup github-copilot --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".github"
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
			}

			outputFile := fmt.Sprintf("%s/copilot-instructions.md", dirPath)

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
			markdownContent := fmt.Sprintf(`# GitHub Copilot Instructions

## ctx Documentation

%s
`, content)

			// Write the file
			if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			fmt.Println("\nThese instructions will be automatically included for all Copilot interactions in this repository.")
			return nil
		},
	}

	copilotCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return copilotCmd
}
