package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewWindsurfCmd creates the setup windsurf subcommand
func NewWindsurfCmd() *cobra.Command {
	var force bool

	windsurfCmd := &cobra.Command{
		Use:   "windsurf",
		Short: "Generate Windsurf rule with ctx documentation",
		Long: `Generate a Windsurf IDE rule file with ctx documentation.

This command creates .windsurf/rules/ctx.md with comprehensive
ctx documentation formatted for Windsurf IDE.

Examples:
  ctx setup windsurf         # Generate .windsurf/rules/ctx.md
  ctx setup windsurf --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".windsurf/rules"
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
			fmt.Println("\nTo enable this rule in Windsurf:")
			fmt.Println("1. Go to Windsurf Settings")
			fmt.Println("2. Find 'AI Rules' section")
			fmt.Println("3. Set 'ctx' rule to 'Always On'")
			return nil
		},
	}

	windsurfCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return windsurfCmd
}