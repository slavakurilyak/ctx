package setup

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// CleanHelpOutput removes ANSI escape codes and trailing whitespace
func CleanHelpOutput(input string) string {
	// Remove all ANSI escape sequences with regex
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	output := re.ReplaceAllString(input, "")

	// Split into lines and remove trailing whitespace (preserve intentional indentation)
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		// Remove trailing regular spaces and tabs only
		lines[i] = strings.TrimRight(line, " \t")
	}

	return strings.Join(lines, "\n")
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
