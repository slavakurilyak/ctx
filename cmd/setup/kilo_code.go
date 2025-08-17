package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewKiloCodeCmd creates the setup kilo-code subcommand
func NewKiloCodeCmd() *cobra.Command {
	var force bool

	kiloCodeCmd := &cobra.Command{
		Use:   "kilo-code",
		Short: "Generate Kilo Code rules with ctx documentation",
		Long: `Generate Kilo Code rule files with ctx documentation.

This command creates .kilocode/rules/ctx.md with comprehensive
ctx documentation formatted for Kilo Code.

Examples:
  ctx setup kilo-code         # Generate .kilocode/rules/ctx.md
  ctx setup kilo-code --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".kilocode/rules"
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
			fmt.Println("\nKilo Code will automatically load these rules from the .kilocode/rules/ directory.")
			return nil
		},
	}

	kiloCodeCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return kiloCodeCmd
}