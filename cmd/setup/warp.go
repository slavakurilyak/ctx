package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewWarpCmd creates the setup warp subcommand
func NewWarpCmd() *cobra.Command {
	var force bool

	warpCmd := &cobra.Command{
		Use:   "warp",
		Short: "Generate WARP.md with ctx documentation",
		Long: `Generate a WARP.md file with ctx documentation for Warp Terminal.

This command creates a WARP.md file in the current directory
with comprehensive ctx documentation formatted as Warp project-scoped rules.

Warp automatically detects and applies WARP.md files in project roots,
making this integration seamless and requiring no additional configuration.

Examples:
  ctx setup warp         # Generate WARP.md
  ctx setup warp --force # Overwrite existing WARP.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := "WARP.md"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			markdownContent, err := GenerateRules("Warp Terminal Project Rules - ctx")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			
			// Provide usage instructions
			fmt.Println("\nWarp Terminal integration:")
			fmt.Println("• WARP.md is automatically detected and applied as project rules")
			fmt.Println("• Rules apply automatically when working in this project directory")
			fmt.Println("• Use `/init` in Warp to view or manage project rules")
			fmt.Println("• View all rules in Warp Drive: Personal > Rules > Project-based")
			fmt.Println("\nThe ctx rules ensure Warp's agents generate token-aware commands for cost-effective interactions.")
			
			return nil
		},
	}

	warpCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return warpCmd
}