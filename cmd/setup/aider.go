package setup

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// NewAiderCmd creates the setup aider subcommand
func NewAiderCmd() *cobra.Command {
	var force bool

	aiderCmd := &cobra.Command{
		Use:   "aider",
		Short: "Generate Aider configuration with ctx documentation",
		Long: `Generate an Aider configuration file with ctx documentation.

This command creates .aider.conf.yml with ctx documentation
embedded as comments for reference during Aider sessions.

Examples:
  ctx setup aider         # Generate .aider.conf.yml
  ctx setup aider --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := ".aider.conf.yml"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Get ctx help output
			content, err := GetCtxHelp()
			if err != nil {
				return err
			}

			// Convert content to YAML comments
			lines := strings.Split(content, "\n")
			var commentedLines []string
			for _, line := range lines {
				commentedLines = append(commentedLines, "# "+line)
			}
			commentedContent := strings.Join(commentedLines, "\n")

			// Format as YAML with embedded documentation
			yamlContent := fmt.Sprintf(`# Aider Configuration with ctx Documentation
#
# ctx Documentation (for reference):
%s
#
# ═══════════════════════════════════════════════════════════════════════════════
#
# You can add project-specific Aider settings below:
# auto-commits: true
# dark-mode: true
# model: claude-opus-4-20250514
# edit-format: diff
#
# For supported models, see:
# https://github.com/Aider-AI/aider/blob/main/aider/models.py

# Example of adding files to always include in context:
# read:
#   - README.md
#   - CONVENTIONS.md
`, commentedContent)

			// Write the file
			if err := os.WriteFile(outputFile, []byte(yamlContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			return nil
		},
	}

	aiderCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return aiderCmd
}