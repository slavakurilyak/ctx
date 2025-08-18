package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewRooCodeCmd creates the setup roo-code subcommand
func NewRooCodeCmd() *cobra.Command {
	var force bool

	rooCodeCmd := &cobra.Command{
		Use:   "roo-code",
		Short: "Generate Roo Code rules with ctx documentation",
		Long: `Generate Roo Code rule files with ctx documentation.

This command creates .roo/rules/ctx.md with comprehensive
ctx documentation formatted for Roo Code.

Examples:
  ctx setup roo-code         # Generate .roo/rules/ctx.md
  ctx setup roo-code --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".roo/rules"
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
			}

			outputFile := fmt.Sprintf("%s/ctx.md", dirPath)

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
			fmt.Println("\nRoo Code will automatically load these rules from the .roo/rules/ directory.")
			return nil
		},
	}

	rooCodeCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return rooCodeCmd
}
