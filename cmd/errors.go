package cmd

import (
	"fmt"
	"strings"
)

// ImproveErrorMessage enhances error messages with helpful suggestions
func ImproveErrorMessage(err error, originalArgs []string) error {
	if err == nil {
		return nil
	}

	errStr := err.Error()

	// Check for common flag parsing errors
	if strings.Contains(errStr, "unknown shorthand flag") ||
		strings.Contains(errStr, "unknown flag") ||
		strings.Contains(errStr, "unknown command") {

		// Build the command string for examples
		cmdStr := strings.Join(originalArgs, " ")

		return fmt.Errorf(`%s

It looks like ctx is trying to parse your command's flags as its own.

Try one of these methods:

  1. Using run subcommand:
     ctx run %s
     ctx --max-tokens 5000 run %s
  
  2. Using quotes (for simple commands):
     ctx "%s"
     ctx --max-tokens 5000 "%s"
  
  3. Using -- separator (POSIX standard):
     ctx -- %s
     ctx --max-tokens 5000 -- %s

Remember: ctx flags (like --max-tokens) must come BEFORE the separator or subcommand.`,
			err.Error(),
			cmdStr,
			cmdStr,
			cmdStr,
			cmdStr,
			cmdStr,
			cmdStr)
	}

	// Return original error if no improvement needed
	return err
}
