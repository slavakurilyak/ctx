package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCrushCmd creates the setup crush subcommand
func NewCrushCmd() *cobra.Command {
	var force bool

	crushCmd := &cobra.Command{
		Use:   "crush",
		Short: "Generate CRUSH.md with ctx documentation",
		Long: `Generate a CRUSH.md file with ctx documentation for Crush Terminal.

This command creates a CRUSH.md file in the current directory
with comprehensive ctx documentation that Crush automatically detects and uses as context.

Crush searches for context files automatically, making this integration
seamless and requiring no additional configuration.

Examples:
  ctx setup crush         # Generate CRUSH.md
  ctx setup crush --force # Overwrite existing CRUSH.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := "CRUSH.md"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			markdownContent, err := GenerateRules("Crush Terminal Context - ctx")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			
			// Provide usage instructions
			fmt.Println("\nCrush Terminal integration:")
			fmt.Println("• CRUSH.md is automatically detected and loaded as context")
			fmt.Println("• No additional configuration required in Crush")
			fmt.Println("• Context applies to all LLM interactions within this project")
			fmt.Println("• Crush also supports variants like crush.md, Crush.md, etc.")
			fmt.Println("\nThe ctx context ensures Crush generates token-aware commands for cost-effective AI interactions.")
			
			return nil
		},
	}

	crushCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return crushCmd
}