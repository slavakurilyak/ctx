package setup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// NewOpenCodeCmd creates the setup opencode subcommand
func NewOpenCodeCmd() *cobra.Command {
	var force bool

	openCodeCmd := &cobra.Command{
		Use:   "opencode",
		Short: "Generate AGENTS.md with ctx documentation",
		Long: `Generate an AGENTS.md file with ctx documentation for OpenCode AI assistant.

This command creates an AGENTS.md file in the current directory
with ctx help output for OpenCode to understand ctx usage.

Examples:
  ctx setup opencode         # Generate AGENTS.md
  ctx setup opencode --force # Overwrite existing AGENTS.md`,
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFile := "AGENTS.md"

			// Check if file exists and force flag is not set
			if _, err := os.Stat(outputFile); err == nil && !force {
				return fmt.Errorf("file %s already exists. Use --force to overwrite", outputFile)
			}

			// Get ctx help output
			cmdCtx := exec.Command("ctx", "-h")
			output, err := cmdCtx.Output()
			if err != nil {
				return fmt.Errorf("failed to get ctx help output: %v", err)
			}

			// Format as markdown with raw ctx help
			content := fmt.Sprintf(`# ctx Documentation

%s
`, strings.TrimSpace(string(output)))

			// Write the file
			if err := os.WriteFile(outputFile, []byte(content), 0644); err != nil {
				return fmt.Errorf("failed to write %s: %v", outputFile, err)
			}

			fmt.Printf("Documentation generated: %s\n", outputFile)
			return nil
		},
	}

	openCodeCmd.Flags().BoolVar(&force, "force", false, "Overwrite existing file")

	return openCodeCmd
}