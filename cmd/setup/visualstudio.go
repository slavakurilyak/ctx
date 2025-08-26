package setup

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewVisualStudioCmd creates the setup visualstudio subcommand
func NewVisualStudioCmd() *cobra.Command {
	var force bool

	visualStudioCmd := &cobra.Command{
		Use:   "visualstudio",
		Short: "Generate Visual Studio 2022 custom instructions with ctx documentation",
		Long: `Generate custom instructions file for Visual Studio 2022 with ctx documentation.

This command creates .github/copilot-instructions.md with comprehensive
ctx documentation formatted for GitHub Copilot Chat in Visual Studio 2022.

The generated file is compatible with:
- Visual Studio 2022 (with GitHub Copilot Chat)
- GitHub.com (Copilot in browser)
- GitHub Copilot CLI

Examples:
  ctx setup visualstudio         # Generate .github/copilot-instructions.md
  ctx setup visualstudio --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create the directory structure
			dirPath := ".github"
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", dirPath, err)
			}

			outputFile := fmt.Sprintf("%s/copilot-instructions.md", dirPath)

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Generate rules using the shared generator
			content, err := GenerateRules("Visual Studio 2022 Custom Instructions")
			if err != nil {
				return fmt.Errorf("failed to generate rules: %v", err)
			}

			// Format as markdown without frontmatter (VS 2022 doesn't use frontmatter)
			markdownContent := fmt.Sprintf(`# GitHub Copilot Custom Instructions

## ctx Documentation

%s

## About These Instructions

These custom instructions are automatically included in all GitHub Copilot Chat interactions within this repository when using:
- Visual Studio 2022 with GitHub Copilot Chat
- GitHub.com (web interface)
- GitHub Copilot CLI

The instructions help GitHub Copilot understand that this project uses ctx for token-aware command execution.
`, content)

			// Write the file
			if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			fmt.Println("\nTo enable custom instructions in Visual Studio 2022:")
			fmt.Println("1. Open Visual Studio 2022")
			fmt.Println("2. Go to Tools > Options > GitHub > Copilot")
			fmt.Println("3. Enable '(Preview) Enable custom instructions to be loaded from .github/copilot-instructions.md files'")
			fmt.Println("4. Restart Visual Studio 2022")
			fmt.Println("\nThe ctx instructions will now be automatically applied to all GitHub Copilot Chat interactions.")
			fmt.Println("\nNote: This file is also compatible with GitHub Copilot on GitHub.com and GitHub Copilot CLI.")
			
			return nil
		},
	}

	visualStudioCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return visualStudioCmd
}