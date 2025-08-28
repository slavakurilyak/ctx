package cmd

import (
	"fmt"

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
- Non-interactive mode: Use --non-interactive to default to AGENTS.md (universal format)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no arguments provided, check for non-interactive mode
			if len(args) == 0 {
				nonInteractive, _ := cmd.Flags().GetBool("non-interactive")
				if nonInteractive {
					// Default to AGENTS.md for universal compatibility
					return setupAgentsMd(force)
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
	setupCmd.Flags().Bool("non-interactive", false, "Run in non-interactive mode (defaults to AGENTS.md)")

	// Add subcommands for each tool
	setupCmd.AddCommand(setup.NewAgentsCmd())  // Generic AGENTS.md (OpenAI standard)
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
	setupCmd.AddCommand(setup.NewVSCodeCmd())
	setupCmd.AddCommand(setup.NewVisualStudioCmd())
	setupCmd.AddCommand(setup.NewAugmentCodeCmd())
	setupCmd.AddCommand(setup.NewAmazonQCmd())
	setupCmd.AddCommand(setup.NewZencoderCmd())
	setupCmd.AddCommand(setup.NewQodoCmd())
	setupCmd.AddCommand(setup.NewWarpCmd())
	setupCmd.AddCommand(setup.NewCrushCmd())

	return setupCmd
}

// setupAgentsMd handles the default AGENTS.md setup
func setupAgentsMd(force bool) error {
	// Reuse the agents command implementation
	agentsCmd := setup.NewAgentsCmd()
	if force {
		agentsCmd.Flags().Set("force", "true")
	}
	return agentsCmd.RunE(agentsCmd, []string{})
}

// showAvailableTools displays available tools for human users
func showAvailableTools() error {
	fmt.Println("Available AI Tools for Setup:")
	fmt.Println()
	fmt.Println("  ctx setup agents        # AGENTS.md (OpenAI standard - works with multiple agents)")
	fmt.Println()
	fmt.Println("  ctx setup claude        # Claude Code/Desktop")
	fmt.Println("  ctx setup cursor        # Cursor IDE")
	fmt.Println("  ctx setup aider         # Aider")
	fmt.Println("  ctx setup windsurf      # Windsurf IDE")
	fmt.Println("  ctx setup jetbrains     # JetBrains AI Assistant")
	fmt.Println("  ctx setup augmentcode   # Augment Code")
	fmt.Println("  ctx setup gemini        # Gemini CLI")
	fmt.Println("  ctx setup zed           # Zed Editor")
	fmt.Println("  ctx setup github-copilot # GitHub Copilot")
	fmt.Println("  ctx setup vscode        # VS Code with GitHub Copilot")
	fmt.Println("  ctx setup visualstudio  # Visual Studio 2022")
	fmt.Println("  ctx setup amazonq       # Amazon Q Developer")
	fmt.Println("  ctx setup zencoder      # Zencoder")
	fmt.Println("  ctx setup qodo          # Qodo Gen")
	fmt.Println("  ctx setup warp          # Warp Terminal")
	fmt.Println("  ctx setup crush         # Crush Terminal")
	fmt.Println("  ctx setup cline         # Cline")
	fmt.Println("  ctx setup roo-code      # Roo Code")
	fmt.Println("  ctx setup kilo-code     # Kilo Code")
	fmt.Println("  ctx setup goose         # Goose AI Agent")
	fmt.Println("  ctx setup trae          # Trae IDE")
	fmt.Println("  ctx setup opencode      # OpenCode")
	fmt.Println()
	fmt.Println("Choose the tool you're using, then run the specific command.")
	fmt.Println("For AI agents: Use 'ctx setup --non-interactive' to default to AGENTS.md.")
	fmt.Println()

	return nil
}
