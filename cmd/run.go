package cmd

import (
	"context"
	"fmt"

	"github.com/slavakurilyak/ctx/internal/app"
	"github.com/spf13/cobra"
)

// NewRunCmd creates the run subcommand for explicit command execution
func NewRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run [command] [args...]",
		Short: "Execute a command with ctx wrapping",
		Long: `Execute a command with ctx wrapping. All arguments after 'run' are passed to the command.

This is an alternative to using the -- separator.

Examples:
  ctx run ls -la
  ctx --max-tokens 5000 run psql -c "SELECT * FROM users"
  ctx --no-tokens run docker ps`,
		Args: cobra.ArbitraryArgs, // Accept any number of args
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get parent command (root command)
			parentCmd := cmd.Parent()
			if parentCmd == nil {
				return fmt.Errorf("internal error: no parent command")
			}
			
			// Check if we have a command to run
			if len(args) == 0 {
				return fmt.Errorf("no command specified after 'run'\n\nUsage: ctx run <command> [args...]\nExample: ctx run ls -la")
			}
			
			// Use the parent's context which has AppContext
			// We need to get it from the parent's context since PersistentPreRunE runs on parent
			appCtx, ok := parentCmd.Context().Value(app.AppContextKey).(*app.AppContext)
			if !ok || appCtx == nil {
				// Try from our own context as fallback
				appCtx, ok = cmd.Context().Value(app.AppContextKey).(*app.AppContext)
				if !ok || appCtx == nil {
					return fmt.Errorf("application context not initialized")
				}
			}
			
			executor := NewCommandExecutor(appCtx)
			
			// Check if streaming is enabled on parent
			isStream, _ := parentCmd.Flags().GetBool("stream")
			if isStream {
				return executor.ExecuteStreamCommand(context.Background(), args)
			}
			
			return executor.ExecuteCommand(context.Background(), args)
		},
	}
}