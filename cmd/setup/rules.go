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
const RulesVersion = "0.1.1"

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
		// Skip ONLY the ctx Directives section (keep everything after INVOCATION)
		if strings.Contains(line, "ctx Directives") {
			skipDirectives = true
			continue
		}
		if skipDirectives {
			if strings.Contains(line, "INVOCATION:") {
				skipDirectives = false
				// Include this line and everything after
			} else {
				continue
			}
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
	return `## ctx Overview

This project uses **ctx** - a universal tool wrapper that provides token awareness for AI agents. 

ctx wraps any CLI tool, shell command, or script to provide:
- Precise token counting (OpenAI, Anthropic, Gemini)
- Execution metadata and telemetry
- Safety controls (token limits, output limits)
- Structured JSON output for LLM consumption

**CRITICAL**: ALWAYS prefix commands with ` + "`ctx`" + ` to enable token-aware execution.

` + helpOutput
}