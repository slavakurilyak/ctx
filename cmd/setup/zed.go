package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewZedCmd creates the setup zed subcommand
func NewZedCmd() *cobra.Command {
	var force bool

	zedCmd := &cobra.Command{
		Use:   "zed",
		Short: "Generate Zed rule with ctx documentation",
		Long: `Generate a Zed editor rule file with ctx documentation.

This command creates a .rules file with comprehensive
ctx documentation formatted for Zed's Agent Panel.

Examples:
  ctx setup zed         # Generate .rules
  ctx setup zed --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := ".rules"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			rulesContent, err := GenerateRules(".rules")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(rulesContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			fmt.Println("\nThese rules will be automatically included in all Agent Panel interactions in Zed.")
			return nil
		},
	}

	zedCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return zedCmd
}
