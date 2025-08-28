package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewQodoCmd creates the setup qodo subcommand
func NewQodoCmd() *cobra.Command {
	var force bool

	qodoCmd := &cobra.Command{
		Use:   "qodo",
		Short: "Generate QODO.md with ctx documentation",
		Long: `Generate a QODO.md file with ctx documentation for Qodo Gen.

This command creates a QODO.md file in the current directory
with comprehensive ctx documentation formatted for Qodo Gen custom instructions.

You can copy the content from QODO.md and paste it into Qodo Gen's Chat Preferences
to enable ctx token-awareness across all interactions.

Examples:
  ctx setup qodo         # Generate QODO.md
  ctx setup qodo --force # Overwrite existing QODO.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := "QODO.md"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			markdownContent, err := GenerateRules("Qodo Gen Custom Instructions - ctx")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			
			// Provide setup instructions
			fmt.Println("\nTo use these instructions in Qodo Gen:")
			fmt.Println("1. Open Qodo Gen Chat by clicking the Qodo Gen icon")
			fmt.Println("2. Click the three dots icon on the top right")
			fmt.Println("3. Choose 'Chat Preferences'")
			fmt.Println("4. Under 'Custom Instructions' → 'All messages', paste the content from QODO.md")
			fmt.Println("5. Save your preferences")
			fmt.Println("\nThe ctx instructions will now apply to all Qodo Gen interactions, ensuring token-aware command execution.")
			
			return nil
		},
	}

	qodoCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return qodoCmd
}