package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewClineCmd creates the setup cline subcommand
func NewClineCmd() *cobra.Command {
	var force bool

	clineCmd := &cobra.Command{
		Use:   "cline",
		Short: "Generate Cline rules with ctx documentation",
		Long: `Generate a Cline rules file with ctx documentation.

This command creates a .clinerules file with comprehensive
ctx documentation formatted for Cline.

Cline uses .clinerules files for system-level guidance and
project-specific context that persists across conversations.

Examples:
  ctx setup cline         # Generate .clinerules
  ctx setup cline --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := ".clinerules"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Get ctx help output
			content, err := GetCtxHelp()
			if err != nil {
				return err
			}

			// Format as markdown for Cline rules
			rulesContent := fmt.Sprintf(`# ctx Documentation

%s
`, content)

			// Write the file
			if err := os.WriteFile(outputFile, []byte(rulesContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			fmt.Println("\nCline will automatically use these rules for all conversations in this project.")
			fmt.Println("For global rules, place this file in ~/Documents/Cline/Rules/")
			return nil
		},
	}

	clineCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return clineCmd
}
