package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/fang"
	"github.com/joho/godotenv"
	"github.com/slavakurilyak/ctx/cmd"
)

// These variables are replaced by ldflags during build
// They are kept here for compatibility but internal/version is the source of truth

func main() {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	// We no longer pre-build the AppContext here.
	// It will be constructed in PersistentPreRunE after flags are parsed.
	rootCmd := cmd.NewRootCmdWithDI()

	// Disable Cobra's default error handling to manage exit codes manually
	rootCmd.SilenceErrors = true

	// fang.Execute handles the CLI execution, including styled errors and panics.
	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		var exitErr *cmd.ExitError
		// Check for our custom exit code error for wrapped commands
		if errors.As(err, &exitErr) {
			// This is an expected failure from the wrapped command.
			// The JSON output has already been printed.
			// We just need to exit with the correct code.
			os.Exit(exitErr.Code) // This will be ExitCodeWrappedCmdError (1).
		}

		// Check if this is a flag parsing error and improve the message
		errStr := err.Error()
		if strings.Contains(errStr, "unknown shorthand flag") ||
			strings.Contains(errStr, "unknown flag") ||
			strings.Contains(errStr, "unknown command") {
			// The error message will be improved by ImproveErrorMessage if called
			// But in case it wasn't, provide helpful guidance here too
			if !strings.Contains(errStr, "Try one of these methods") {
				fmt.Fprintf(os.Stderr, `%v

Hint: Use one of these methods:
  ctx run ls -la
  ctx "ls -la"
  ctx -- ls -la
`, err)
			} else {
				fmt.Fprintf(os.Stderr, "%v\n", err)
			}
		} else {
			// For all other errors, treat them as internal application errors.
			// This would include tokenizer initialization, config loading, etc.
			fmt.Fprintf(os.Stderr, "ctx application error: %v\n", err)
		}
		os.Exit(cmd.ExitCodeAppError)
	}
}
