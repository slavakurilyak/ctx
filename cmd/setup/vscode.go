package setup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// NewVSCodeCmd creates the setup vscode subcommand
func NewVSCodeCmd() *cobra.Command {
	var force bool
	var userProfile bool

	vscodeCmd := &cobra.Command{
		Use:   "vscode",
		Short: "Generate VS Code custom instructions file with ctx documentation",
		Long: `Generate a VS Code custom instructions file with ctx documentation.

This command creates .github/instructions/ctx.instructions.md with comprehensive
ctx documentation formatted for VS Code's GitHub Copilot Chat.

The file includes frontmatter configuration to automatically apply the instructions
to all files in your workspace when using GitHub Copilot in VS Code.

Examples:
  ctx setup vscode                # Generate .github/instructions/ctx.instructions.md
  ctx setup vscode --force         # Overwrite existing file
  ctx setup vscode --user-profile  # Install to user profile instead of workspace`,
		RunE: func(cmd *cobra.Command, args []string) error {
			var dirPath string
			var outputFile string

			if userProfile {
				// User profile location (cross-workspace)
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get user home directory: %v", err)
				}
				dirPath = filepath.Join(homeDir, ".vscode", "instructions")
				outputFile = filepath.Join(dirPath, "ctx.instructions.md")
			} else {
				// Workspace location (default)
				dirPath = ".github/instructions"
				outputFile = filepath.Join(dirPath, "ctx.instructions.md")
			}

			// Create the directory structure
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
			}

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			baseContent, err := GenerateRules("VS Code Custom Instructions")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Add VS Code-specific frontmatter
			frontmatter := `---
applyTo: "**"
description: "ctx - Universal tool wrapper for token-aware command execution"
---

`
			fullContent := frontmatter + baseContent

			// Write the file
			if err := os.WriteFile(outputFile, []byte(fullContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			
			// Provide configuration instructions
			fmt.Println("\nTo enable custom instructions in VS Code:")
			fmt.Println("1. Open VS Code Settings (Cmd/Ctrl+,)")
			fmt.Println("2. Search for 'github.copilot.chat.codeGeneration.useInstructionFiles'")
			fmt.Println("3. Enable the setting")
			fmt.Println("4. Restart VS Code or reload the window")
			fmt.Println("\nThe ctx instructions will now be automatically applied to all GitHub Copilot interactions.")
			
			if userProfile {
				fmt.Println("\nNote: Instructions installed to user profile will apply to all workspaces.")
			} else {
				fmt.Println("\nNote: Instructions installed to workspace will only apply to this project.")
				fmt.Println("\nThis file is also compatible with GitHub Copilot on GitHub.com when stored in .github/copilot-instructions.md")
			}
			
			return nil
		},
	}

	vscodeCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")
	vscodeCmd.Flags().BoolVar(&userProfile, "user-profile", false, "Install to user profile instead of workspace")

	return vscodeCmd
}