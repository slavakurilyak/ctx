package cmd

import (
	"fmt"
	"os"

	"github.com/slavakurilyak/ctx/cmd/setup"
	"github.com/spf13/cobra"
)

// NewSetupCmd creates the setup command
func NewSetupCmd() *cobra.Command {
	var force bool

	setupCmd := &cobra.Command{
		Use:   "setup [tool]",
		Short: "Set up coding agents and assistants with ctx documentation",
		Long: `Generate ctx documentation for various coding agents and assistants.

Available tools: claude, cursor, aider, windsurf, jetbrains, gemini, etc.

When no tool is specified:
- Interactive mode: Shows available tools and prompts for selection
- Non-interactive mode: Use --non-interactive to default to Claude Code`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no arguments provided, check for non-interactive mode
			if len(args) == 0 {
				nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
				if nonInteractive {
					// Default to Claude Code for AI agents
					return setupClaudeCode(force)
				} else {
					// Show available tools for human users
					return showAvailableTools()
				}
			}
			// Otherwise show help for available subcommands
			return cmd.Help()
		},
	}

	setupCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")
	setupCmd.Flags().Bool("non-interactive", false, "Run in non-interactive mode (defaults to Claude Code)")

	// Add subcommands for each tool
	setupCmd.AddCommand(setup.NewClaudeCmd())
	setupCmd.AddCommand(setup.NewCursorCmd())
	setupCmd.AddCommand(setup.NewWindsurfCmd())
	setupCmd.AddCommand(setup.NewJetBrainsCmd())
	setupCmd.AddCommand(setup.NewAiderCmd())
	setupCmd.AddCommand(setup.NewZedCmd())
	setupCmd.AddCommand(setup.NewGitHubCopilotCmd())
	setupCmd.AddCommand(setup.NewRooCodeCmd())
	setupCmd.AddCommand(setup.NewKiloCodeCmd())
	setupCmd.AddCommand(setup.NewGeminiCmd())
	setupCmd.AddCommand(setup.NewClineCmd())
	setupCmd.AddCommand(setup.NewGooseCmd())
	setupCmd.AddCommand(setup.NewTraeCmd())
	setupCmd.AddCommand(setup.NewOpenCodeCmd())

	return setupCmd
}

// setupClaudeCode handles the default Claude Code setup
func setupClaudeCode(force bool) error {
	outputFile := "CLAUDE.md"

	// Check if file exists and force flag is not set
	if _, err := os.Stat(outputFile); err == nil && !force {
		return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
	}

	// Get ctx help output
	content, err := setup.GetCtxHelp()
	if err != nil {
		return err
	}

	// Format as markdown
	markdownContent := fmt.Sprintf(`# ctx Documentation

%s
`, content)

	// Write the file
	if err := os.WriteFile(outputFile, []byte(markdownContent), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %v", outputFile, err)
	}

	fmt.Printf("Documentation generated: %s\n", outputFile)
	return nil
}

// showAvailableTools displays available tools for human users
func showAvailableTools() error {
	fmt.Println("Available AI Tools for Setup:")
	fmt.Println()
	fmt.Println("  ctx setup claude        # Claude Code/Desktop")
	fmt.Println("  ctx setup cursor        # Cursor IDE")
	fmt.Println("  ctx setup aider         # Aider")
	fmt.Println("  ctx setup windsurf      # Windsurf IDE")
	fmt.Println("  ctx setup jetbrains     # JetBrains AI Assistant")
	fmt.Println("  ctx setup gemini        # Gemini CLI")
	fmt.Println("  ctx setup zed           # Zed Editor")
	fmt.Println("  ctx setup github-copilot # GitHub Copilot")
	fmt.Println("  ctx setup cline         # Cline")
	fmt.Println("  ctx setup roo-code      # Roo Code")
	fmt.Println("  ctx setup kilo-code     # Kilo Code")
	fmt.Println("  ctx setup goose         # Goose AI Agent")
	fmt.Println("  ctx setup trae          # Trae IDE")
	fmt.Println("  ctx setup opencode      # OpenCode")
	fmt.Println()
	fmt.Println("Choose the tool you're using, then run the specific command.")
	fmt.Println("For AI agents: Use 'ctx setup --non-interactive' to default to Claude Code.")
	fmt.Println()

	return nil
}
