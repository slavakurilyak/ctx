package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewJetBrainsCmd creates the setup jetbrains subcommand
func NewJetBrainsCmd() *cobra.Command {
	var force bool

	jetbrainsCmd := &cobra.Command{
		Use:   "jetbrains",
		Short: "Generate JetBrains AI Assistant rule with ctx documentation",
		Long: `Generate a JetBrains AI Assistant rule file with ctx documentation.

This command creates .aiassistant/rules/ctx.md with comprehensive
ctx documentation formatted for JetBrains AI Assistant.

Examples:
  ctx setup jetbrains         # Generate .aiassistant/rules/ctx.md
  ctx setup jetbrains --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".aiassistant/rules"
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
			fmt.Println("\nTo enable this rule in JetBrains AI Assistant:")
			fmt.Println("1. Go to Settings → AI Assistant")
			fmt.Println("2. Find 'Custom Rules' section")
			fmt.Println("3. Set 'ctx' rule type to 'Always'")
			return nil
		},
	}

	jetbrainsCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return jetbrainsCmd
}