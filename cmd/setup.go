package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// cleanHelpOutput removes excessive whitespace from help text
func cleanHelpOutput(input string) string {
	lines := strings.Split(input, "\n")
	var cleaned []string
	emptyCount := 0
	
	for _, line := range lines {
		// Only trim trailing spaces, keep leading spaces as-is
		line = strings.TrimRight(line, " ")
		
		// Handle multiple empty lines - allow max 1 blank line
		if line == "" {
			emptyCount++
			if emptyCount <= 1 {
				cleaned = append(cleaned, line)
			}
		} else {
			emptyCount = 0
			cleaned = append(cleaned, line)
		}
	}
	
	return strings.Join(cleaned, "\n")
}

// NewSetupCmd creates the setup command
func NewSetupCmd() *cobra.Command {
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: "Set up coding agents and assistants",
		Long:  "Set up coding agents and assistants with ctx documentation and best practices",
	}

	setupCmd.AddCommand(newSetupClaudeCmd())
	setupCmd.AddCommand(newSetupCursorCmd())
	setupCmd.AddCommand(newSetupWindsurfCmd())
	setupCmd.AddCommand(newSetupJetBrainsCmd())
	setupCmd.AddCommand(newSetupAiderCmd())
	setupCmd.AddCommand(newSetupZedCmd())
	setupCmd.AddCommand(newSetupGitHubCopilotCmd())
	setupCmd.AddCommand(newSetupAllCmd())

	return setupCmd
}

// newSetupClaudeCmd creates the setup claude subcommand
func newSetupClaudeCmd() *cobra.Command {
	var force bool
	var output string

	claudeCmd := &cobra.Command{
		Use:   "claude [filename]",
		Short: "Set up Claude AI with ctx documentation",
		Long: `Set up Claude AI assistant with ctx documentation and best practices.

Examples:
  ctx setup claude                    # Generate CLAUDE.md
  ctx setup claude --output README.md # Generate README.md
  ctx setup claude --force            # Overwrite existing file`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Determine output filename
			filename := "CLAUDE.md"
			if output != "" {
				filename = output
			} else if len(args) > 0 {
				filename = args[0]
			}

			// Check if file exists and --force not provided
			if !force {
				if _, err := os.Stat(filename); err == nil {
					return fmt.Errorf("file %s already exists. Use --force to overwrite", filename)
				}
			}

			// Get ctx help output from current binary
			// Try to use the current binary if it exists, otherwise fall back to PATH
			ctxBinary := "./ctx"
			if _, err := os.Stat(ctxBinary); os.IsNotExist(err) {
				ctxBinary = "ctx"
			}
			cmdCtx := exec.Command(ctxBinary, "-h")
			output, err := cmdCtx.Output()
			if err != nil {
				return fmt.Errorf("failed to get ctx help output: %v", err)
			}

			// Clean up the output (remove extra whitespace)
			content := strings.TrimSpace(string(output))
			content = cleanHelpOutput(content)
			
			// Create markdown content without code block
			markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)

			// Write to file
			err = os.WriteFile(filename, []byte(markdownContent), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", filename, err)
			}

			fmt.Printf("Documentation generated: %s\n", filename)
			return nil
		},
	}

	claudeCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")
	claudeCmd.Flags().StringVarP(&output, "output", "o", "", "Output filename")

	return claudeCmd
}

// newSetupCursorCmd creates the setup cursor subcommand
func newSetupCursorCmd() *cobra.Command {
	var force bool

	cursorCmd := &cobra.Command{
		Use:   "cursor",
		Short: "Set up Cursor IDE with ctx documentation",
		Long: `Set up Cursor IDE AI assistant with ctx documentation and best practices.

This creates a Cursor rule that will always be applied to provide
ctx documentation to the AI assistant.

Examples:
  ctx setup cursor         # Generate .cursor/rules/ctx.mdc
  ctx setup cursor --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create .cursor/rules directory if it doesn't exist
			rulesDir := ".cursor/rules"
			if err := os.MkdirAll(rulesDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", rulesDir, err)
			}

			filename := fmt.Sprintf("%s/ctx.mdc", rulesDir)

			// Check if file exists and --force not provided
			if !force {
				if _, err := os.Stat(filename); err == nil {
					return fmt.Errorf("file %s already exists. Use --force to overwrite", filename)
				}
			}

			// Get ctx help output from current binary
			// Try to use the current binary if it exists, otherwise fall back to PATH
			ctxBinary := "./ctx"
			if _, err := os.Stat(ctxBinary); os.IsNotExist(err) {
				ctxBinary = "ctx"
			}
			cmdCtx := exec.Command(ctxBinary, "-h")
			output, err := cmdCtx.Output()
			if err != nil {
				return fmt.Errorf("failed to get ctx help output: %v", err)
			}

			// Clean up the output
			content := strings.TrimSpace(string(output))
			content = cleanHelpOutput(content)
			
			// Create MDC content with metadata header
			mdcContent := fmt.Sprintf(`---
alwaysApply: true
---

# ctx Documentation

%s
`, content)

			// Write to file
			err = os.WriteFile(filename, []byte(mdcContent), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", filename, err)
			}

			fmt.Printf("Cursor rule generated: %s\n", filename)
			return nil
		},
	}

	cursorCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return cursorCmd
}

// newSetupWindsurfCmd creates the setup windsurf subcommand
func newSetupWindsurfCmd() *cobra.Command {
	var force bool

	windsurfCmd := &cobra.Command{
		Use:   "windsurf",
		Short: "Set up Windsurf IDE with ctx documentation",
		Long: `Set up Windsurf IDE AI assistant with ctx documentation and best practices.

This creates a Windsurf rule that can be set to "Always On" to provide
ctx documentation to the AI assistant.

Examples:
  ctx setup windsurf         # Generate .windsurf/rules/ctx.md
  ctx setup windsurf --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create .windsurf/rules directory if it doesn't exist
			rulesDir := ".windsurf/rules"
			if err := os.MkdirAll(rulesDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", rulesDir, err)
			}

			filename := fmt.Sprintf("%s/ctx.md", rulesDir)

			// Check if file exists and --force not provided
			if !force {
				if _, err := os.Stat(filename); err == nil {
					return fmt.Errorf("file %s already exists. Use --force to overwrite", filename)
				}
			}

			// Get ctx help output from current binary
			// Try to use the current binary if it exists, otherwise fall back to PATH
			ctxBinary := "./ctx"
			if _, err := os.Stat(ctxBinary); os.IsNotExist(err) {
				ctxBinary = "ctx"
			}
			cmdCtx := exec.Command(ctxBinary, "-h")
			output, err := cmdCtx.Output()
			if err != nil {
				return fmt.Errorf("failed to get ctx help output: %v", err)
			}

			// Clean up the output
			content := strings.TrimSpace(string(output))
			content = cleanHelpOutput(content)
			
			// Create markdown content
			markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)

			// Write to file
			err = os.WriteFile(filename, []byte(markdownContent), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", filename, err)
			}

			fmt.Printf("Windsurf rule generated: %s\n", filename)
			return nil
		},
	}

	windsurfCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return windsurfCmd
}

// newSetupJetBrainsCmd creates the setup jetbrains subcommand
func newSetupJetBrainsCmd() *cobra.Command {
	var force bool

	jetbrainsCmd := &cobra.Command{
		Use:   "jetbrains",
		Short: "Set up JetBrains AI Assistant with ctx documentation",
		Long: `Set up JetBrains AI Assistant with ctx documentation and best practices.

This creates a JetBrains AI Assistant rule that will be automatically applied
to provide ctx documentation to the AI assistant.

Examples:
  ctx setup jetbrains         # Generate .aiassistant/rules/ctx.md
  ctx setup jetbrains --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create .aiassistant/rules directory if it doesn't exist
			rulesDir := ".aiassistant/rules"
			if err := os.MkdirAll(rulesDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", rulesDir, err)
			}

			filename := fmt.Sprintf("%s/ctx.md", rulesDir)

			// Check if file exists and --force not provided
			if !force {
				if _, err := os.Stat(filename); err == nil {
					return fmt.Errorf("file %s already exists. Use --force to overwrite", filename)
				}
			}

			// Get ctx help output from current binary
			// Try to use the current binary if it exists, otherwise fall back to PATH
			ctxBinary := "./ctx"
			if _, err := os.Stat(ctxBinary); os.IsNotExist(err) {
				ctxBinary = "ctx"
			}
			cmdCtx := exec.Command(ctxBinary, "-h")
			output, err := cmdCtx.Output()
			if err != nil {
				return fmt.Errorf("failed to get ctx help output: %v", err)
			}

			// Clean up the output
			content := strings.TrimSpace(string(output))
			content = cleanHelpOutput(content)
			
			// Create markdown content
			markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)

			// Write to file
			err = os.WriteFile(filename, []byte(markdownContent), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", filename, err)
			}

			fmt.Printf("JetBrains AI Assistant rule generated: %s\n", filename)
			return nil
		},
	}

	jetbrainsCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return jetbrainsCmd
}

// newSetupZedCmd creates the setup zed subcommand
func newSetupZedCmd() *cobra.Command {
	var force bool

	zedCmd := &cobra.Command{
		Use:   "zed",
		Short: "Set up Zed editor with ctx documentation",
		Long: `Set up Zed editor AI assistant with ctx documentation and best practices.

This creates a Zed rules file that will be automatically included
to provide ctx documentation to the AI assistant.

Examples:
  ctx setup zed         # Generate .rules
  ctx setup zed --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := ".rules"

			// Check if file exists and --force not provided
			if !force {
				if _, err := os.Stat(filename); err == nil {
					return fmt.Errorf("file %s already exists. Use --force to overwrite", filename)
				}
			}

			// Get ctx help output from current binary
			// Try to use the current binary if it exists, otherwise fall back to PATH
			ctxBinary := "./ctx"
			if _, err := os.Stat(ctxBinary); os.IsNotExist(err) {
				ctxBinary = "ctx"
			}
			cmdCtx := exec.Command(ctxBinary, "-h")
			output, err := cmdCtx.Output()
			if err != nil {
				return fmt.Errorf("failed to get ctx help output: %v", err)
			}

			// Clean up the output
			content := strings.TrimSpace(string(output))
			content = cleanHelpOutput(content)
			
			// Create markdown content
			markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)

			// Write to file
			err = os.WriteFile(filename, []byte(markdownContent), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", filename, err)
			}

			fmt.Printf("Zed rules generated: %s\n", filename)
			return nil
		},
	}

	zedCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return zedCmd
}

// newSetupGitHubCopilotCmd creates the setup github-copilot subcommand
func newSetupGitHubCopilotCmd() *cobra.Command {
	var force bool

	copilotCmd := &cobra.Command{
		Use:   "github-copilot",
		Short: "Set up GitHub Copilot with ctx documentation",
		Long: `Set up GitHub Copilot with ctx documentation and best practices.

This creates GitHub Copilot instructions that will be automatically included
to provide ctx documentation to GitHub Copilot across the repository.

Examples:
  ctx setup github-copilot         # Generate .github/copilot-instructions.md
  ctx setup github-copilot --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create .github directory if it doesn't exist
			githubDir := ".github"
			if err := os.MkdirAll(githubDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", githubDir, err)
			}

			filename := fmt.Sprintf("%s/copilot-instructions.md", githubDir)

			// Check if file exists and --force not provided
			if !force {
				if _, err := os.Stat(filename); err == nil {
					return fmt.Errorf("file %s already exists. Use --force to overwrite", filename)
				}
			}

			// Get ctx help output from current binary
			// Try to use the current binary if it exists, otherwise fall back to PATH
			ctxBinary := "./ctx"
			if _, err := os.Stat(ctxBinary); os.IsNotExist(err) {
				ctxBinary = "ctx"
			}
			cmdCtx := exec.Command(ctxBinary, "-h")
			output, err := cmdCtx.Output()
			if err != nil {
				return fmt.Errorf("failed to get ctx help output: %v", err)
			}

			// Clean up the output
			content := strings.TrimSpace(string(output))
			content = cleanHelpOutput(content)
			
			// Create markdown content
			markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)

			// Write to file
			err = os.WriteFile(filename, []byte(markdownContent), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", filename, err)
			}

			fmt.Printf("GitHub Copilot instructions generated: %s\n", filename)
			return nil
		},
	}

	copilotCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return copilotCmd
}

// newSetupAiderCmd creates the setup aider subcommand
func newSetupAiderCmd() *cobra.Command {
	var force bool

	aiderCmd := &cobra.Command{
		Use:   "aider",
		Short: "Set up Aider with ctx documentation",
		Long: `Set up Aider AI assistant with ctx documentation and best practices.

This creates a self-contained Aider configuration file that includes
all ctx documentation as comments, ensuring Aider always has context
about ctx usage patterns and best practices.

Examples:
  ctx setup aider         # Generate .aider.conf.yml
  ctx setup aider --force # Overwrite existing file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := ".aider.conf.yml"

			// Check if file exists and --force not provided
			if !force {
				if _, err := os.Stat(filename); err == nil {
					return fmt.Errorf("file %s already exists. Use --force to overwrite", filename)
				}
			}

			// Get ctx help output from current binary
			// Try to use the current binary if it exists, otherwise fall back to PATH
			ctxBinary := "./ctx"
			if _, err := os.Stat(ctxBinary); os.IsNotExist(err) {
				ctxBinary = "ctx"
			}
			cmdCtx := exec.Command(ctxBinary, "-h")
			output, err := cmdCtx.Output()
			if err != nil {
				return fmt.Errorf("failed to get ctx help output: %v", err)
			}

			// Clean up the output
			content := strings.TrimSpace(string(output))
			content = cleanHelpOutput(content)
			
			// Convert content to YAML comments (prefix each line with #)
			lines := strings.Split(content, "\n")
			var commentedLines []string
			for _, line := range lines {
				if line == "" {
					commentedLines = append(commentedLines, "#")
				} else {
					commentedLines = append(commentedLines, "# " + line)
				}
			}
			documentationComments := strings.Join(commentedLines, "\n")

			// Create YAML content with embedded documentation
			yamlContent := fmt.Sprintf(`# ============================================================================
# ctx Documentation - The Context Engine for AI
# ============================================================================
#
# ALWAYS wrap commands with ctx to track token usage and prevent cost explosions.
# ctx provides precise token counts, execution metadata, and safety controls.
#
%s
#
# ============================================================================
# Aider Configuration
# ============================================================================

# This configuration file includes ctx documentation to ensure proper usage.
# When working with commands, remember to use ctx to wrap them for token awareness.

# You can add project-specific Aider settings below:
# auto-commits: true
# dark-mode: true
# model: gpt-4o
# edit-format: diff

# Example of adding files to always include in context:
# read:
#   - README.md
#   - CONVENTIONS.md
`, documentationComments)

			// Write to file
			err = os.WriteFile(filename, []byte(yamlContent), 0644)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %v", filename, err)
			}

			fmt.Printf("Aider configuration generated: %s\n", filename)
			fmt.Println("ctx documentation is embedded in the configuration for reference")
			return nil
		},
	}

	aiderCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return aiderCmd
}

// newSetupAllCmd creates the setup all subcommand
func newSetupAllCmd() *cobra.Command {
	var force bool

	allCmd := &cobra.Command{
		Use:   "all",
		Short: "Set up all coding agents and assistants",
		Long: `Set up all supported coding agents and assistants with ctx documentation.

This command sets up:
  - CLAUDE.md
  - .cursor/rules/ctx.mdc
  - .windsurf/rules/ctx.md
  - .aiassistant/rules/ctx.md
  - .aider.conf.yml
  - .rules (for Zed)
  - .github/copilot-instructions.md

Examples:
  ctx setup all         # Set up all coding assistants
  ctx setup all --force # Overwrite existing files`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get ctx help output once
			cmdCtx := exec.Command("ctx", "-h")
			output, err := cmdCtx.Output()
			if err != nil {
				return fmt.Errorf("failed to get ctx help output: %v", err)
			}

			// Clean up the output
			content := strings.TrimSpace(string(output))
			content = cleanHelpOutput(content)

			// Generate CLAUDE.md
			claudeFile := "CLAUDE.md"
			if !force {
				if _, err := os.Stat(claudeFile); err == nil {
					fmt.Printf("Skipping %s (already exists, use --force to overwrite)\n", claudeFile)
				} else {
					markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
					if err := os.WriteFile(claudeFile, []byte(markdownContent), 0644); err != nil {
						return fmt.Errorf("failed to write %s: %v", claudeFile, err)
					}
					fmt.Printf("Documentation generated: %s\n", claudeFile)
				}
			} else {
				markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
				if err := os.WriteFile(claudeFile, []byte(markdownContent), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %v", claudeFile, err)
				}
				fmt.Printf("Documentation generated: %s\n", claudeFile)
			}

			// Generate Cursor rule
			cursorDir := ".cursor/rules"
			if err := os.MkdirAll(cursorDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", cursorDir, err)
			}
			cursorFile := fmt.Sprintf("%s/ctx.mdc", cursorDir)
			
			if !force {
				if _, err := os.Stat(cursorFile); err == nil {
					fmt.Printf("Skipping %s (already exists, use --force to overwrite)\n", cursorFile)
				} else {
					mdcContent := fmt.Sprintf(`---
alwaysApply: true
---

# ctx Documentation

%s
`, content)
					if err := os.WriteFile(cursorFile, []byte(mdcContent), 0644); err != nil {
						return fmt.Errorf("failed to write %s: %v", cursorFile, err)
					}
					fmt.Printf("Cursor rule generated: %s\n", cursorFile)
				}
			} else {
				mdcContent := fmt.Sprintf(`---
alwaysApply: true
---

# ctx Documentation

%s
`, content)
				if err := os.WriteFile(cursorFile, []byte(mdcContent), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %v", cursorFile, err)
				}
				fmt.Printf("Cursor rule generated: %s\n", cursorFile)
			}

			// Generate Windsurf rule
			windsurfDir := ".windsurf/rules"
			if err := os.MkdirAll(windsurfDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", windsurfDir, err)
			}
			windsurfFile := fmt.Sprintf("%s/ctx.md", windsurfDir)
			
			if !force {
				if _, err := os.Stat(windsurfFile); err == nil {
					fmt.Printf("Skipping %s (already exists, use --force to overwrite)\n", windsurfFile)
				} else {
					markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
					if err := os.WriteFile(windsurfFile, []byte(markdownContent), 0644); err != nil {
						return fmt.Errorf("failed to write %s: %v", windsurfFile, err)
					}
					fmt.Printf("Windsurf rule generated: %s\n", windsurfFile)
				}
			} else {
				markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
				if err := os.WriteFile(windsurfFile, []byte(markdownContent), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %v", windsurfFile, err)
				}
				fmt.Printf("Windsurf rule generated: %s\n", windsurfFile)
			}

			// Generate JetBrains AI Assistant rule
			jetbrainsDir := ".aiassistant/rules"
			if err := os.MkdirAll(jetbrainsDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", jetbrainsDir, err)
			}
			jetbrainsFile := fmt.Sprintf("%s/ctx.md", jetbrainsDir)
			
			if !force {
				if _, err := os.Stat(jetbrainsFile); err == nil {
					fmt.Printf("Skipping %s (already exists, use --force to overwrite)\n", jetbrainsFile)
				} else {
					markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
					if err := os.WriteFile(jetbrainsFile, []byte(markdownContent), 0644); err != nil {
						return fmt.Errorf("failed to write %s: %v", jetbrainsFile, err)
					}
					fmt.Printf("JetBrains AI Assistant rule generated: %s\n", jetbrainsFile)
				}
			} else {
				markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
				if err := os.WriteFile(jetbrainsFile, []byte(markdownContent), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %v", jetbrainsFile, err)
				}
				fmt.Printf("JetBrains AI Assistant rule generated: %s\n", jetbrainsFile)
			}

			// Generate Zed rules
			zedFile := ".rules"
			
			if !force {
				if _, err := os.Stat(zedFile); err == nil {
					fmt.Printf("Skipping %s (already exists, use --force to overwrite)\n", zedFile)
				} else {
					markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
					if err := os.WriteFile(zedFile, []byte(markdownContent), 0644); err != nil {
						return fmt.Errorf("failed to write %s: %v", zedFile, err)
					}
					fmt.Printf("Zed rules generated: %s\n", zedFile)
				}
			} else {
				markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
				if err := os.WriteFile(zedFile, []byte(markdownContent), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %v", zedFile, err)
				}
				fmt.Printf("Zed rules generated: %s\n", zedFile)
			}

			// Generate Aider configuration with embedded documentation
			aiderFile := ".aider.conf.yml"
			
			if !force {
				if _, err := os.Stat(aiderFile); err == nil {
					fmt.Printf("Skipping %s (already exists, use --force to overwrite)\n", aiderFile)
				} else {
					// Convert content to YAML comments
					lines := strings.Split(content, "\n")
					var commentedLines []string
					for _, line := range lines {
						if line == "" {
							commentedLines = append(commentedLines, "#")
						} else {
							commentedLines = append(commentedLines, "# " + line)
						}
					}
					documentationComments := strings.Join(commentedLines, "\n")

					yamlContent := fmt.Sprintf(`# ============================================================================
# ctx Documentation - The Context Engine for AI
# ============================================================================
#
# ALWAYS wrap commands with ctx to track token usage and prevent cost explosions.
# ctx provides precise token counts, execution metadata, and safety controls.
#
%s
#
# ============================================================================
# Aider Configuration
# ============================================================================

# This configuration file includes ctx documentation to ensure proper usage.
# When working with commands, remember to use ctx to wrap them for token awareness.

# You can add project-specific Aider settings below:
# auto-commits: true
# dark-mode: true
# model: gpt-4o
# edit-format: diff

# Example of adding files to always include in context:
# read:
#   - README.md
#   - CONVENTIONS.md
`, documentationComments)
					
					if err := os.WriteFile(aiderFile, []byte(yamlContent), 0644); err != nil {
						return fmt.Errorf("failed to write %s: %v", aiderFile, err)
					}
					fmt.Printf("Aider configuration generated: %s\n", aiderFile)
				}
			} else {
				// Convert content to YAML comments
				lines := strings.Split(content, "\n")
				var commentedLines []string
				for _, line := range lines {
					if line == "" {
						commentedLines = append(commentedLines, "#")
					} else {
						commentedLines = append(commentedLines, "# " + line)
					}
				}
				documentationComments := strings.Join(commentedLines, "\n")

				yamlContent := fmt.Sprintf(`# ============================================================================
# ctx Documentation - The Context Engine for AI
# ============================================================================
#
# ALWAYS wrap commands with ctx to track token usage and prevent cost explosions.
# ctx provides precise token counts, execution metadata, and safety controls.
#
%s
#
# ============================================================================
# Aider Configuration
# ============================================================================

# This configuration file includes ctx documentation to ensure proper usage.
# When working with commands, remember to use ctx to wrap them for token awareness.

# You can add project-specific Aider settings below:
# auto-commits: true
# dark-mode: true
# model: gpt-4o
# edit-format: diff

# Example of adding files to always include in context:
# read:
#   - README.md
#   - CONVENTIONS.md
`, documentationComments)
				
				if err := os.WriteFile(aiderFile, []byte(yamlContent), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %v", aiderFile, err)
				}
				fmt.Printf("Aider configuration generated: %s\n", aiderFile)
			}

			// Generate GitHub Copilot instructions
			githubDir := ".github"
			if err := os.MkdirAll(githubDir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", githubDir, err)
			}
			copilotFile := fmt.Sprintf("%s/copilot-instructions.md", githubDir)
			
			if !force {
				if _, err := os.Stat(copilotFile); err == nil {
					fmt.Printf("Skipping %s (already exists, use --force to overwrite)\n", copilotFile)
				} else {
					markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
					if err := os.WriteFile(copilotFile, []byte(markdownContent), 0644); err != nil {
						return fmt.Errorf("failed to write %s: %v", copilotFile, err)
					}
					fmt.Printf("GitHub Copilot instructions generated: %s\n", copilotFile)
				}
			} else {
				markdownContent := fmt.Sprintf("# ctx Documentation\n\n%s\n", content)
				if err := os.WriteFile(copilotFile, []byte(markdownContent), 0644); err != nil {
					return fmt.Errorf("failed to write %s: %v", copilotFile, err)
				}
				fmt.Printf("GitHub Copilot instructions generated: %s\n", copilotFile)
			}

			return nil
		},
	}

	allCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")

	return allCmd
}