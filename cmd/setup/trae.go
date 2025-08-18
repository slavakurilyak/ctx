package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewTraeCmd creates the setup trae subcommand
func NewTraeCmd() *cobra.Command {
	var force bool

	traeCmd := &cobra.Command{
		Use:   "trae",
		Short: "Generate Trae rules with ctx documentation",
		Long: `Generate a Trae IDE rules file with ctx documentation.

This command creates .trae/rules/project_rules.md with comprehensive
ctx documentation formatted for Trae IDE's AI Rules system.

Examples:
  ctx setup trae         # Generate .trae/rules/project_rules.md
  ctx setup trae --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".trae/rules"
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
			}

			outputFile := fmt.Sprintf("%s/project_rules.md", dirPath)

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
			return nil
		},
	}

	traeCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return traeCmd
}
