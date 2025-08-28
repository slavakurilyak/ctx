package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewAmazonQCmd creates the setup amazonq subcommand
func NewAmazonQCmd() *cobra.Command {
	var force bool

	amazonqCmd := &cobra.Command{
		Use:   "amazonq",
		Short: "Generate Amazon Q Developer project rules with ctx documentation",
		Long: `Generate project rules for Amazon Q Developer with ctx documentation.

This command creates .amazonq/rules/ctx.md with comprehensive ctx documentation
formatted as Amazon Q Developer project rules.

Amazon Q will automatically use these rules as context when developers chat
within your project, ensuring consistency with coding standards and best practices.

Examples:
  ctx setup amazonq         # Generate .amazonq/rules/ctx.md
  ctx setup amazonq --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".amazonq/rules"
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
			}

			outputFile := filepath.Join(dirPath, "ctx.md")

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			content, err := GenerateRules("Amazon Q Developer Project Rules - ctx")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			
			// Provide configuration instructions
			fmt.Println("\nTo use project rules in Amazon Q Developer:")
			fmt.Println("1. The rules are automatically detected in the .amazonq/rules/ folder")
			fmt.Println("2. Amazon Q will use these rules as context in chat interactions")
			fmt.Println("3. In the chat interface, click the 'Rules' button to:")
			fmt.Println("   - View all available rules")
			fmt.Println("   - Toggle rules on/off for the current session")
			fmt.Println("   - Create additional rules")
			fmt.Println("\nThe ctx rules will ensure Amazon Q generates code consistent with your project's token awareness requirements.")
			
			return nil
		},
	}

	amazonqCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return amazonqCmd
}