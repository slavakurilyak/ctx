package setup

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/slavakurilyak/ctx/internal/models"
	"github.com/slavakurilyak/ctx/internal/version"
)

// RulesVersion defines the version of the rules template
// This is incremented when the structure or content of generated rules changes
const RulesVersion = "0.1.0"

// GenerateRules generates standardized rules/instructions for any agent
func GenerateRules(agentName string) (string, error) {
	// Get ctx help output
	cmdCtx := exec.Command("ctx", "-h")
	output, err := cmdCtx.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get ctx help output: %v", err)
	}

	helpOutput := strings.TrimSpace(string(output))
	cleanedHelp := cleanHelpOutput(helpOutput)

	// Generate version header
	header := fmt.Sprintf(`# %s
<!-- ctx-version: %s -->
<!-- rules-version: %s -->
<!-- schema-version: %s -->
<!-- generated: %s -->

`, agentName, version.GetVersion(), RulesVersion, models.CurrentSchemaVersion, time.Now().Format("2006-01-02"))

	// Generate content with all sections
	content := generateContent(cleanedHelp)

	return header + content, nil
}

// cleanHelpOutput removes duplicated sections from help output
func cleanHelpOutput(helpOutput string) string {
	lines := strings.Split(helpOutput, "\n")
	var cleanedHelp []string
	skipDirectives := false
	skipRecommended := false
	skipFilterTypes := false

	for _, line := range lines {
		// Skip the ctx Directives section (duplicate of Project Overview)
		if strings.Contains(line, "ctx Directives") {
			skipDirectives = true
			continue
		}
		if skipDirectives && strings.Contains(line, "USAGE:") {
			skipDirectives = false
		}
		if skipDirectives {
			continue
		}

		// Skip RECOMMENDED USAGE section (duplicate of Essential ctx Patterns)
		if strings.Contains(line, "RECOMMENDED USAGE:") {
			skipRecommended = true
			continue
		}
		if skipRecommended && strings.Contains(line, "FILTER TYPES:") {
			skipRecommended = false
			skipFilterTypes = true
			continue
		}
		if skipRecommended {
			continue
		}

		// Skip FILTER TYPES section (duplicate of Setting Limits and Filters)
		if skipFilterTypes && strings.Contains(line, "════") {
			skipFilterTypes = false
			continue
		}
		if skipFilterTypes {
			continue
		}

		cleanedHelp = append(cleanedHelp, line)
	}

	return strings.Join(cleanedHelp, "\n")
}

// generateContent creates the main content sections
func generateContent(helpOutput string) string {
	return `## Project Overview

This project uses **ctx** - a universal tool wrapper that provides token awareness for AI agents. 

ctx wraps any CLI tool, shell command, or script to provide:
- Precise token counting (OpenAI, Anthropic, Gemini)
- Execution metadata and telemetry
- Safety controls (token limits, output limits)
- Structured JSON output for LLM consumption

**CRITICAL**: ALWAYS prefix commands with ` + "`ctx`" + ` to enable token-aware execution.

## Setup Commands

- Install ctx: See https://github.com/slavakurilyak/ctx for installation
- Configure for your agent: ` + "`ctx setup [agent-name]`" + `
- Verify installation: ` + "`ctx version`" + `

## Essential ctx Patterns

### Always Wrap Commands

**ALWAYS** wrap commands with ctx for token awareness:
- **Good:** ` + "`ctx npm test`" + `
- **Good:** ` + "`ctx docker ps`" + `
- **Good:** ` + "`ctx git status`" + `
- **Bad:** ` + "`npm test`" + ` (never run directly)

### Probe Before Acting

Before running expensive commands, probe their token cost:
` + "```bash" + `
# Check token cost first
ctx psql -c "SELECT * FROM users" | jq '.tokens'

# If too expensive, refine the query
ctx --max-tokens 5000 psql -c "SELECT id, name FROM users LIMIT 100"
` + "```" + `

### Common Command Examples

All commands should be wrapped with ctx:
- Run tests: ` + "`ctx npm test`" + ` or ` + "`ctx go test ./...`" + `
- Build project: ` + "`ctx npm run build`" + ` or ` + "`ctx go build`" + `
- Start dev server: ` + "`ctx npm run dev`" + `
- Lint code: ` + "`ctx npm run lint`" + `
- Database queries: ` + "`ctx psql -c \"SELECT ...\"`" + `
- Docker operations: ` + "`ctx docker ps`" + `, ` + "`ctx docker logs app`" + `

### Setting Limits and Filters

**Token Limits:**
` + "```bash" + `
ctx --max-tokens 5000 command        # Limit output tokens
ctx --max-lines 100 command          # Limit output lines
ctx --max-output-bytes 1048576       # Limit output size (1MB)
ctx --max-pipeline-stages 5 command  # Limit pipeline stages
` + "```" + `

**Filtering Strategies:**
` + "```bash" + `
# Command-level filtering (more efficient)
ctx docker logs app --tail 100
ctx psql -c "SELECT * FROM users LIMIT 10"

# ctx-level filtering (safety net)
ctx --max-tokens 1000 docker logs app
` + "```" + `

## ctx Output Schema

ctx returns structured JSON with schema version ` + models.CurrentSchemaVersion + `:

` + "```json" + `
{
  "tokens": 296,
  "output": "command output here",
  "input": "psql -U user -l",
  "metadata": {
    "exit_code": 0,
    "success": true,
    "duration": 29,
    "bytes": 1092,
    "timestamp": "2025-08-20T00:14:30-04:00",
    "directory": "/current/working/directory",
    "user": "username",
    "host": "hostname.local",
    "session_id": "fbed68b6-fb56-4f65-9da0-fb071bf8b63c",
    "limits": {
      "max_tokens": 5000,
      "actual_lines": 42,
      "limit_reached": "token_limit"
    }
  },
  "telemetry": {
    "trace_id": "abc123",
    "span_id": "def456",
    "trace_flags": "01"
  },
  "schema_version": "` + models.CurrentSchemaVersion + `"
}
` + "```" + `

### Schema Fields

- **tokens**: Token count for the output
- **output**: The actual command output
- **input**: The command that was executed
- **metadata**: Execution details and context
  - **exit_code**: Process exit code
  - **success**: Boolean (true if exit_code is 0)
  - **duration**: Execution time in milliseconds
  - **bytes**: Output size in bytes
  - **timestamp**: RFC3339 formatted timestamp
  - **directory**: Working directory
  - **user**: Username who ran the command
  - **host**: Hostname
  - **session_id**: Unique session identifier
  - **limits**: Information about applied limits (optional)
- **telemetry**: OpenTelemetry trace information (optional)
- **schema_version**: Version of the output schema

## Security Considerations

- Never expose secrets in command outputs
- Use ` + "`--private`" + ` flag for sensitive operations
- Token counts help prevent accidental data exposure
- Set appropriate limits for production environments

## Additional Resources

- ctx repository: https://github.com/slavakurilyak/ctx
- Report issues: https://github.com/slavakurilyak/ctx/issues

## ctx Command Reference

` + helpOutput
}