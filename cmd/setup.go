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

When no tool is specified, defaults to Claude Code (generates CLAUDE.md).`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no arguments provided, default to Claude Code setup
			if len(args) == 0 {
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
			// Otherwise show help for available subcommands
			return cmd.Help()
		},
	}
	
	setupCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing files")

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

	return setupCmd
}