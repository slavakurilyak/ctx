package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewAugmentCodeCmd creates the setup augmentcode subcommand
func NewAugmentCodeCmd() *cobra.Command {
	var force bool
	var legacy bool

	augmentCodeCmd := &cobra.Command{
		Use:   "augmentcode",
		Short: "Generate Augment Code rules with ctx documentation",
		Long: `Generate Augment Code rules with ctx documentation.

This command creates rules that Augment Code will automatically include in every
Agent and Chat session for token-aware command execution.

By default, creates .augment/rules/ctx.md (recommended format).
Use --legacy flag to create .augment-guidelines file instead.

Rules are automatically imported by Augment Code and marked as "Always" rules,
ensuring ctx instructions are included in every interaction.

Examples:
  ctx setup augmentcode           # Generate .augment/rules/ctx.md
  ctx setup augmentcode --force   # Overwrite existing files
  ctx setup augmentcode --legacy  # Use legacy .augment-guidelines format`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var outputFile string
			var dirPath string

			if legacy {
				// Legacy workspace guidelines
				outputFile = ".augment-guidelines"
			} else {
				// Modern rules format
				dirPath = filepath.Join(".augment", "rules")
				if err := os.MkdirAll(dirPath, 0755); err != nil {
					return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
				}
				outputFile = filepath.Join(dirPath, "ctx.md")
			}

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			baseContent, err := GenerateRules("Augment Code Rules")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Format content based on the file type
			var content string
			if legacy {
				// Legacy format - simple guidelines
				content = fmt.Sprintf(`# Augment Code Workspace Guidelines

## ctx Documentation

%s

## About These Guidelines

These workspace guidelines are automatically included in all Augment Code Agent and Chat sessions
within this repository. They help Augment Code understand that this project uses ctx for
token-aware command execution.

Note: This is the legacy .augment-guidelines format. Consider using the newer .augment/rules/
directory structure for better organization and flexibility.
`, baseContent)
			} else {
				// Modern rules format with metadata
				content = fmt.Sprintf(`---
type: always
description: ctx - Universal tool wrapper for token-aware command execution
---

# ctx Rules for Augment Code

%s

## About This Rule

This rule is automatically included in all Augment Code Agent and Chat sessions.
It ensures that Augment Code understands how to use ctx for token-aware command
execution in this project.

Rule Type: **Always** - This rule is included in every user message to ensure
consistent use of ctx across all interactions.
`, baseContent)
			}

			// Write the file
			if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			
			if legacy {
				fmt.Println("\nLegacy workspace guidelines created.")
				fmt.Println("Augment Code will automatically import these guidelines for all Agent and Chat sessions.")
				fmt.Println("\nConsider migrating to the newer .augment/rules/ format for better flexibility.")
			} else {
				fmt.Println("\nAugment Code will automatically import this rule.")
				fmt.Println("The rule is marked as 'always' and will be included in every Agent and Chat interaction.")
				fmt.Println("\nTo verify the rule is imported:")
				fmt.Println("1. Open Augment Code in VS Code")
				fmt.Println("2. Click the hamburger menu (⋯) > Settings")
				fmt.Println("3. Select 'User Guidelines and Rules' from the left menu")
				fmt.Println("4. Look for 'ctx.md' under imported rules")
			}
			
			return nil
		},
	}

	augmentCodeCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")
	augmentCodeCmd.Flags().BoolVar(&legacy, "legacy", false, "Use legacy .augment-guidelines format")

	return augmentCodeCmd
}