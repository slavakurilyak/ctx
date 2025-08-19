package setup

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

// CleanHelpOutput removes ANSI escape codes and formats with explicit code blocks
func CleanHelpOutput(input string) string {
	// Remove all ANSI escape sequences with regex
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	output := re.ReplaceAllString(input, "")

	// Split into lines and process with code block grouping
	lines := strings.Split(output, "\n")
	var result []string
	var inCodeBlock bool
	
	i := 0
	for i < len(lines) {
		line := strings.TrimRight(lines[i], " \t")
		
		// Handle empty lines
		if line == "" {
			if inCodeBlock {
				result = append(result, "")
			} else {
				result = append(result, "")
			}
			i++
			continue
		}
		
		// Clean leading whitespace
		cleanLine := strings.TrimLeft(line, " \t\u00a0")
		
		// Special handling for ENVIRONMENT section
		if strings.HasPrefix(cleanLine, "ENVIRONMENT:") {
			if inCodeBlock {
				result = append(result, "```")
				result = append(result, "")
				inCodeBlock = false
			}
			result = append(result, formatHeader(cleanLine))
			result = append(result, "") // Add line break after Environment header
			i++
			
			// Collect and format environment variables
			var currentEnvVar string
			var currentDesc string
			for i < len(lines) {
				envLine := strings.TrimRight(lines[i], " \t")
				cleanEnvLine := strings.TrimLeft(envLine, " \t\u00a0")
				
				// Stop at next section header
				if cleanEnvLine != "" && isHeaderLine(cleanEnvLine) && !strings.HasPrefix(cleanEnvLine, "ENVIRONMENT:") {
					break
				}
				
				// Process environment variable lines
				if cleanEnvLine != "" && !strings.HasPrefix(cleanEnvLine, "CLI flags") {
					if isEnvironmentVariable(cleanEnvLine) {
						// If we have a previous env var, add it to results
						if currentEnvVar != "" {
							result = append(result, fmt.Sprintf("- `%s` - %s", currentEnvVar, currentDesc))
						}
						// Start tracking new env var
						currentEnvVar = cleanEnvLine
						currentDesc = ""
					} else if currentEnvVar != "" {
						// This is a description for the current env var
						currentDesc = cleanEnvLine
					}
				} else if strings.HasPrefix(cleanEnvLine, "CLI flags") {
					// Add the CLI flags note
					result = append(result, cleanEnvLine)
					result = append(result, "")
				}
				
				i++
			}
			
			// Add the last env var if there is one
			if currentEnvVar != "" {
				result = append(result, fmt.Sprintf("- `%s` - %s", currentEnvVar, currentDesc))
				result = append(result, "")
			}
			continue
		}
		
		shouldBeInCode := shouldBeInCodeBlock(cleanLine)
		
		if shouldBeInCode {
			// Start code block if not already in one
			if !inCodeBlock {
				result = append(result, "")
				result = append(result, "```")
				inCodeBlock = true
			}
			
			// Add the line to code block
			result = append(result, cleanLine)
			
			// Look ahead to see if we should continue the code block
			nextShouldBeCode := false
			if i+1 < len(lines) {
				nextLine := strings.TrimRight(lines[i+1], " \t")
				if nextLine == "" {
					// Skip empty line and check the line after
					if i+2 < len(lines) {
						nextNextLine := strings.TrimLeft(strings.TrimRight(lines[i+2], " \t"), " \t\u00a0")
						nextShouldBeCode = shouldBeInCodeBlock(nextNextLine)
					}
				} else {
					nextCleanLine := strings.TrimLeft(nextLine, " \t\u00a0")
					nextShouldBeCode = shouldBeInCodeBlock(nextCleanLine)
				}
			}
			
			// If next line shouldn't be in code, close the block
			if !nextShouldBeCode {
				result = append(result, "```")
				result = append(result, "")
				inCodeBlock = false
			}
			
		} else {
			// Handle non-code content
			if inCodeBlock {
				result = append(result, "```")
				result = append(result, "")
				inCodeBlock = false
			}
			
			if isHeaderLine(cleanLine) {
				formatted := formatHeader(cleanLine)
				// Add line break before certain headers
				if strings.HasPrefix(formatted, "## Usage") {
					result = append(result, "")
				}
				result = append(result, formatted)
				// Add line break after certain headers
				if strings.HasPrefix(formatted, "## ctx Directives") ||
					strings.HasPrefix(formatted, "## Environment") ||
					strings.HasPrefix(formatted, "## Recommended Usage") ||
					strings.HasPrefix(formatted, "## Filter Types") ||
					strings.HasPrefix(formatted, "## Commands") ||
					strings.HasPrefix(formatted, "## Flags") {
					result = append(result, "")
				}
			} else if isWorkflowHeader(cleanLine) {
				// Format workflow headers as H3 with line break before
				result = append(result, "")
				result = append(result, formatWorkflowHeader(cleanLine))
			} else if isCommandLine(cleanLine) {
				// Format command lines with proper spacing
				formatted := formatCommandLine(cleanLine)
				result = append(result, formatted)
			} else if isFlagLine(cleanLine) {
				formatted := formatFlagLine(cleanLine)
				result = append(result, formatted...)
			} else {
				if strings.TrimSpace(cleanLine) != "" {
					// Special formatting for filter type bullet points
					if strings.HasPrefix(cleanLine, "• ") {
						result = append(result, formatFilterBullet(cleanLine))
					} else if strings.Contains(cleanLine, "════") {
						// Add line break before separator line
						result = append(result, "")
						result = append(result, cleanLine)
					} else if strings.HasPrefix(cleanLine, "ctx") && strings.Contains(cleanLine, "[command]") {
						// Wrap usage line in code block
						result = append(result, "`" + cleanLine + "`")
					} else {
						result = append(result, cleanLine)
					}
				}
			}
		}
		
		i++
	}
	
	// Close any open code block at the end
	if inCodeBlock {
		result = append(result, "```")
	}

	return strings.Join(result, "\n")
}

// shouldBeInCodeBlock determines if a line should be inside a code block
func shouldBeInCodeBlock(line string) bool {
	// Don't treat bullet points as code - they should be formatted with inline code
	if strings.HasPrefix(line, "• ") {
		return false
	}
	
	// Commands starting with ctx
	if strings.HasPrefix(line, "ctx ") {
		return true
	}
	
	// Comments starting with #
	if strings.HasPrefix(line, "# ") {
		return true
	}
	
	// Lines that look like command examples or shell output
	if strings.Contains(line, " | jq ") || strings.Contains(line, "psql -c") {
		return true
	}
	
	// Don't treat individual env vars as code - they'll be grouped together
	return false
}

// isHeaderLine checks if a line is a section header
func isHeaderLine(line string) bool {
	headers := []string{"USAGE:", "EXAMPLES:", "ENVIRONMENT:", "RECOMMENDED USAGE:", 
					   "FILTER TYPES:", "COMMANDS", "FLAGS", "ctx Directives", "USAGE"}
	for _, header := range headers {
		if strings.HasPrefix(line, header) {
			return true
		}
	}
	return false
}

// formatHeader converts a header line to markdown format
func formatHeader(line string) string {
	// Map of headers to their markdown equivalents
	headerMap := map[string]string{
		"USAGE:": "## Usage",
		"USAGE": "## Usage",
		"EXAMPLES:": "## Examples",
		"ENVIRONMENT:": "## Environment",
		"RECOMMENDED USAGE:": "## Recommended Usage",
		"FILTER TYPES:": "## Filter Types",
		"COMMANDS": "## Commands",
		"FLAGS": "## Flags",
	}
	
	// Check if the line matches any header
	for oldHeader, newHeader := range headerMap {
		if strings.HasPrefix(line, oldHeader) {
			return newHeader
		}
	}
	
	// Keep ctx Directives as is (already properly formatted)
	if strings.Contains(line, "ctx Directives") {
		return line
	}
	
	return line
}

// isWorkflowHeader checks if a line is a workflow header (PROBE, PROBE-ACT, etc.)
func isWorkflowHeader(line string) bool {
	workflows := []string{"PROBE", "PROBE-ACT", "PROBE-FILTER-ACT"}
	for _, workflow := range workflows {
		if line == workflow {
			return true
		}
	}
	return false
}

// formatWorkflowHeader converts workflow headers to markdown H3
func formatWorkflowHeader(line string) string {
	workflowMap := map[string]string{
		"PROBE": "### Probe",
		"PROBE-ACT": "### Probe-Act",
		"PROBE-FILTER-ACT": "### Probe-Filter-Act",
	}
	
	if formatted, ok := workflowMap[line]; ok {
		return formatted
	}
	return line
}

// formatFilterBullet formats filter type bullet points consistently with command/flag style
func formatFilterBullet(line string) string {
	// Remove the bullet point prefix
	line = strings.TrimPrefix(line, "• ")
	
	// Handle "ctx filters (safety):" line
	if strings.Contains(line, "ctx filters (safety):") {
		// Format as: - ctx filters (safety) - list of flags
		return "- ctx filters (safety) - `--max-tokens`, `--max-lines`, `--max-output-bytes`, `--timeout`"
	}
	
	// Handle "Command filters (efficiency):" line
	if strings.Contains(line, "Command filters (efficiency):") {
		// Format as: - Command filters (efficiency) - list of filters
		return "- Command filters (efficiency) - `LIMIT`, `WHERE`, `--tail`, `--since`, `head`, `grep`"
	}
	
	// Handle "Best:" line
	if strings.Contains(line, "Best:") {
		// Find the command part and wrap it in backticks
		if idx := strings.Index(line, "ctx "); idx != -1 {
			command := line[idx:]
			return "- Best practice - `" + command + "`"
		}
	}
	
	return "- " + line
}

// isEnvironmentVariable checks if a line is an environment variable name
func isEnvironmentVariable(line string) bool {
	// Environment variables typically start with CTX_ or OTEL_
	return strings.HasPrefix(line, "CTX_") || strings.HasPrefix(line, "OTEL_")
}

// isCommandLine checks if a line is a command description
func isCommandLine(line string) bool {
	// Commands are lines that start with a word followed by spaces and description
	// They appear in the COMMANDS section
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	
	// Check if it looks like a command line (word followed by spaces)
	parts := strings.Fields(trimmed)
	if len(parts) < 2 {
		return false
	}
	
	// Commands don't start with dashes and have descriptions
	if strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "•") {
		return false
	}
	
	// Check if first word is a known command
	commands := []string{"account", "completion", "config", "help", "login", "logout", 
		"run", "setup", "telemetry", "update", "version"}
	for _, cmd := range commands {
		if parts[0] == cmd {
			return true
		}
	}
	
	return false
}

// formatCommandLine formats a command line with proper markdown
func formatCommandLine(line string) string {
	// Find the command name and description
	trimmed := strings.TrimSpace(line)
	
	// Split at multiple spaces to separate command from description
	// Commands have format: "command [args]                    Description"
	parts := regexp.MustCompile(`\s{2,}`).Split(trimmed, 2)
	
	if len(parts) == 2 {
		command := parts[0]
		description := parts[1]
		// Format as: - `command` - Description
		return fmt.Sprintf("- `%s` - %s", command, description)
	}
	
	return "- `" + trimmed + "`"
}

// isFlagLine checks if a line is a flag description
func isFlagLine(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "-") || strings.HasPrefix(trimmed, "--")
}

// formatFlagLine formats a flag line with proper markdown
func formatFlagLine(line string) []string {
	trimmed := strings.TrimSpace(line)
	
	// Split at multiple spaces to separate flag from description
	parts := regexp.MustCompile(`\s{2,}`).Split(trimmed, 2)
	
	if len(parts) == 2 {
		flag := parts[0]
		description := parts[1]
		// Format as: - `flag` - Description
		return []string{fmt.Sprintf("- `%s` - %s", flag, description)}
	}
	
	// If no clear separation, format the whole line
	return []string{"- `" + trimmed + "`"}
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
