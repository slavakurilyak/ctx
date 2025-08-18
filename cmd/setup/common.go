package setup

import (
	"fmt"
	"os/exec"
	"strings"
)

// CleanHelpOutput removes ANSI escape codes
func CleanHelpOutput(input string) string {
	// Remove ANSI escape codes
	output := strings.ReplaceAll(input, "\x1b[0m", "")
	output = strings.ReplaceAll(output, "\x1b[31m", "")
	output = strings.ReplaceAll(output, "\x1b[32m", "")
	output = strings.ReplaceAll(output, "\x1b[33m", "")
	output = strings.ReplaceAll(output, "\x1b[34m", "")
	output = strings.ReplaceAll(output, "\x1b[35m", "")
	output = strings.ReplaceAll(output, "\x1b[36m", "")
	output = strings.ReplaceAll(output, "\x1b[37m", "")
	output = strings.ReplaceAll(output, "\x1b[90m", "")
	output = strings.ReplaceAll(output, "\x1b[1m", "")
	output = strings.ReplaceAll(output, "\x1b[2m", "")
	output = strings.ReplaceAll(output, "\x1b[3m", "")
	output = strings.ReplaceAll(output, "\x1b[4m", "")

	return output
}

// GetCtxHelp gets the ctx help output and cleans it
func GetCtxHelp() (string, error) {
	// Get ctx help output
	cmdCtx := exec.Command("ctx", "-h")
	output, err := cmdCtx.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get ctx help output: %v", err)
	}

	// Clean up the output
	content := strings.TrimSpace(string(output))
	content = CleanHelpOutput(content)

	return content, nil
}
