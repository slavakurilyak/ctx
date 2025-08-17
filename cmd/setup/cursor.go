package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCursorCmd creates the setup cursor subcommand
func NewCursorCmd() *cobra.Command {
	var force bool

	cursorCmd := &cobra.Command{
		Use:   "cursor",
		Short: "Generate Cursor rule with ctx documentation",
		Long: `Generate a Cursor IDE rule file with ctx documentation.

This command creates .cursor/rules/ctx.mdc with comprehensive
ctx documentation formatted for Cursor IDE.

Examples:
  ctx setup cursor         # Generate .cursor/rules/ctx.mdc
  ctx setup cursor --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".cursor/rules"
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
			}

			outputFile := fmt.Sprintf("%s/ctx.mdc", dirPath)

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Get ctx help output
			content, err := GetCtxHelp()
			if err != nil {
				return err
			}

			// Format as MDC (Markdown with metadata)
			mdcContent := fmt.Sprintf(`---
tags: [ctx, ai-tools, command-wrapper, token-efficiency]
alwaysApply: true
---

# ctx Documentation

%s
`, content)

			// Write the file
			if err := os.WriteFile(outputFile, []byte(mdcContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			return nil
		},
	}

	cursorCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return cursorCmd
}