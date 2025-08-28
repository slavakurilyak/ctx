package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewZencoderCmd creates the setup zencoder subcommand
func NewZencoderCmd() *cobra.Command {
	var force bool

	zencoderCmd := &cobra.Command{
		Use:   "zencoder",
		Short: "Generate Zencoder Zen Rules with ctx documentation",
		Long: `Generate Zen Rules for Zencoder with ctx documentation.

This command creates .zencoder/rules/ctx.md with comprehensive ctx documentation
formatted as Zencoder Zen Rules with proper frontmatter.

The rule will be set to always apply to ensure ctx token-awareness is available
across all Zencoder interactions.

Examples:
  ctx setup zencoder         # Generate .zencoder/rules/ctx.md
  ctx setup zencoder --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".zencoder/rules"
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
			}

			outputFile := filepath.Join(dirPath, "ctx.md")

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			baseContent, err := GenerateRules("Zencoder Zen Rules - ctx")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Add Zencoder-specific frontmatter
			frontmatter := `---
description: "ctx - Universal tool wrapper for token-aware command execution"
globs: ["**/*.sh", "**/*.py", "**/*.js", "**/*.go", "**/*.ts", "Makefile", "package.json", "Dockerfile"]
alwaysApply: true
---

`
			fullContent := frontmatter + baseContent

			// Write the file
			if err := os.WriteFile(outputFile, []byte(fullContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			
			// Provide configuration instructions
			fmt.Println("\nTo use Zen Rules in Zencoder:")
			fmt.Println("1. The rule is automatically detected in the .zencoder/rules/ folder")
			fmt.Println("2. With alwaysApply: true, ctx rules will be included in every request")
			fmt.Println("3. You can also manually reference the rule by @mentioning it in chat")
			fmt.Println("4. Type @ and select 'Zen Rules' to see all available rules")
			fmt.Println("\nThe ctx rules ensure Zencoder generates token-aware commands for cost-effective AI interactions.")
			
			return nil
		},
	}

	zencoderCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return zencoderCmd
}