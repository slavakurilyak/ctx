package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/slavakurilyak/ctx/internal/app"
	"github.com/slavakurilyak/ctx/internal/enricher"
	"github.com/slavakurilyak/ctx/internal/executor"
	"github.com/slavakurilyak/ctx/internal/models"
)

const (
	// ExitCodeSuccess indicates ctx and the wrapped command both succeeded.
	ExitCodeSuccess = 0
	// ExitCodeWrappedCmdError indicates ctx ran successfully, but the wrapped command failed.
	ExitCodeWrappedCmdError = 1
	// ExitCodeAppError indicates a failure within ctx itself (e.g., config, tokenizer).
	ExitCodeAppError = 2
	// ExitCodeUsageError indicates a user error, like a bad flag. Cobra handles some of this.
	ExitCodeUsageError = 3
)

// ExitError is used to pass a specific exit code up to main
type ExitError struct {
	Code int
}

func (e *ExitError) Error() string {
	return fmt.Sprintf("command exited with code %d", e.Code)
}

// TokenLimitExceededError is returned when the command's output exceeds the token limit
type TokenLimitExceededError struct {
	Limit  int64
	Actual int
}

func (e *TokenLimitExceededError) Error() string {
	return fmt.Sprintf("token limit of %d exceeded (actual: %d)", e.Limit, e.Actual)
}

// LineLimitExceededError is returned when the line limit is exceeded
type LineLimitExceededError struct {
	Limit int64
}

func (e *LineLimitExceededError) Error() string {
	return fmt.Sprintf("line limit of %d lines exceeded", e.Limit)
}

// OutputLimitExceededError is returned when the byte limit is exceeded
type OutputLimitExceededError struct {
	Limit  int64
	Actual int64
}

func (e *OutputLimitExceededError) Error() string {
	return fmt.Sprintf("output limit of %d bytes exceeded", e.Limit)
}

// CommandExecutor handles command execution with injected dependencies
type CommandExecutor struct {
	enricher *enricher.Enricher
	appCtx   *app.AppContext
}

// NewCommandExecutor creates a new command executor with the given app context
func NewCommandExecutor(appCtx *app.AppContext) *CommandExecutor {
	// Get tokenizer
	tok, _ := appCtx.GetTokenizer() // Ignore error, will work without tokenizer
	
	// Create enricher with dependencies
	enr := enricher.NewEnricher(tok, appCtx.History, appCtx.Telemetry, appCtx.Config)
	
	return &CommandExecutor{
		enricher: enr,
		appCtx:   appCtx,
	}
}

// ExecuteCommand executes a command with the given arguments
func (ce *CommandExecutor) ExecuteCommand(ctx context.Context, args []string) error {
	isPipeline, stages := ce.parseArguments(args)
	
	// Check pipeline stage limit if configured
	if maxStages := ce.appCtx.Config.Limits.MaxPipelineStages; isPipeline && maxStages != nil && len(stages) > *maxStages {
		return fmt.Errorf("pipeline stage limit exceeded: %d stages found, limit is %d", len(stages), *maxStages)
	}
	
	if isPipeline {
		return ce.executePipeline(ctx, stages)
	}
	
	return ce.executeSingleCommand(ctx, args)
}

// parseArguments parses command arguments to detect pipeline mode
func (ce *CommandExecutor) parseArguments(args []string) (bool, [][]string) {
	var stages [][]string
	var currentStage []string
	
	for _, arg := range args {
		if arg == "|" {
			if len(currentStage) > 0 {
				stages = append(stages, currentStage)
				currentStage = []string{}
			}
		} else {
			currentStage = append(currentStage, arg)
		}
	}
	
	if len(currentStage) > 0 {
		stages = append(stages, currentStage)
	}
	
	return len(stages) > 1, stages
}

// executeSingleCommand executes a single command
func (ce *CommandExecutor) executeSingleCommand(ctx context.Context, args []string) error {
	command := strings.Join(args, " ")
	
	// Apply timeout from configuration if set
	if ce.appCtx.Config.DefaultTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ce.appCtx.Config.DefaultTimeout)
		defer cancel()
	}
	
	result, err := executor.ExecuteCommand(ctx, command)
	if err != nil {
		output := models.NewOutput(command, []byte(err.Error()), 1, 0)
		ce.outputResult(output) // Still output the JSON
		return &ExitError{Code: ExitCodeWrappedCmdError}
	}
	
	output, err := ce.enricher.EnrichOutput(ctx, result)
	if err != nil {
		return fmt.Errorf("failed to enrich output: %w", err)
	}
	
	// Post-execution token check
	if ce.appCtx.Config.MaxTokens > 0 && int64(output.Tokens) > ce.appCtx.Config.MaxTokens {
		limitErr := &TokenLimitExceededError{
			Limit:  ce.appCtx.Config.MaxTokens,
			Actual: output.Tokens,
		}
		output.Metadata.Error = limitErr.Error()
		output.Metadata.Success = false
		output.Metadata.FailureReason = "token_limit_exceeded"
		_ = ce.outputResult(output)
		return &ExitError{Code: ExitCodeWrappedCmdError}
	}
	
	err = ce.outputResult(output)
	if err != nil {
		return err
	}
	
	// If the command ran but failed, return the exit code
	if result.ExitCode != 0 {
		return &ExitError{Code: ExitCodeWrappedCmdError}
	}
	
	return nil
}

// executePipeline executes a pipeline of commands
func (ce *CommandExecutor) executePipeline(ctx context.Context, stages [][]string) error {
	pipelineCmd := ce.buildPipelineCommand(stages)
	
	// Apply timeout from configuration if set
	if ce.appCtx.Config.DefaultTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ce.appCtx.Config.DefaultTimeout)
		defer cancel()
	}
	
	result, err := executor.ExecuteCommand(ctx, pipelineCmd)
	if err != nil {
		output := models.NewOutput(pipelineCmd, []byte(err.Error()), 1, 0)
		ce.outputResult(output) // Still output the JSON
		return &ExitError{Code: ExitCodeWrappedCmdError}
	}
	
	output, err := ce.enricher.EnrichOutput(ctx, result)
	if err != nil {
		return fmt.Errorf("failed to enrich pipeline output: %w", err)
	}
	
	// Post-execution token check for pipeline
	if ce.appCtx.Config.MaxTokens > 0 && int64(output.Tokens) > ce.appCtx.Config.MaxTokens {
		limitErr := &TokenLimitExceededError{
			Limit:  ce.appCtx.Config.MaxTokens,
			Actual: output.Tokens,
		}
		output.Metadata.Error = limitErr.Error()
		output.Metadata.Success = false
		output.Metadata.FailureReason = "token_limit_exceeded"
		_ = ce.outputResult(output)
		return &ExitError{Code: ExitCodeWrappedCmdError}
	}
	
	err = ce.outputResult(output)
	if err != nil {
		return err
	}
	
	// If the command ran but failed, return the exit code
	if result.ExitCode != 0 {
		return &ExitError{Code: ExitCodeWrappedCmdError}
	}
	
	return nil
}

// buildPipelineCommand builds a pipeline command from stages
func (ce *CommandExecutor) buildPipelineCommand(stages [][]string) string {
	var commands []string
	for _, stage := range stages {
		commands = append(commands, strings.Join(stage, " "))
	}
	return strings.Join(commands, " | ")
}

// outputResult outputs the result as JSON or pretty format
func (ce *CommandExecutor) outputResult(output *models.Output) error {
	// Check if pretty output is requested
	if ce.appCtx.Config.PrettyOutput {
		return ce.outputPretty(output)
	}
	
	// Default: JSON output
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal output: %w", err)
	}
	
	fmt.Println(string(data))
	return nil
}

// formatBytes formats bytes into human-readable format
func formatBytes(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d bytes", bytes)
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	} else {
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1024*1024))
	}
}

// countLines counts the number of lines in the output
func countLines(s string) int {
	if s == "" {
		return 0
	}
	lines := strings.Count(s, "\n")
	// If the string doesn't end with a newline, it still counts as a line
	if !strings.HasSuffix(s, "\n") {
		lines++
	}
	return lines
}

// getCommandChain extracts command names from the input and formats them as a chain
func getCommandChain(input string) string {
	// Split by pipe to detect pipeline
	parts := strings.Split(input, "|")
	
	if len(parts) == 1 {
		// Single command - extract executable name
		fields := strings.Fields(strings.TrimSpace(parts[0]))
		if len(fields) > 0 {
			// Get base name (handle paths like /usr/bin/psql)
			return filepath.Base(fields[0])
		}
		return "[empty]"
	}
	
	// Pipeline - extract each command's executable
	var commands []string
	for _, part := range parts {
		fields := strings.Fields(strings.TrimSpace(part))
		if len(fields) > 0 {
			commands = append(commands, filepath.Base(fields[0]))
		}
	}
	
	// Join with arrow separator
	if len(commands) > 4 {
		// Truncate long pipelines
		return strings.Join(commands[:3], " → ") + " → ..."
	}
	return strings.Join(commands, " → ")
}

// outputPretty outputs the result in a pretty, terminal-friendly format
func (ce *CommandExecutor) outputPretty(output *models.Output) error {
	// Calculate line count
	lineCount := countLines(output.Output)
	
	// Get command chain for display
	cmdChain := getCommandChain(output.Input)
	
	// Status line with complete metadata (reordered: tokens, time, lines, size)
	if output.Metadata.Success {
		fmt.Printf("ctx → %s | ✓ Success | %d tokens | %d ms | %d lines | %s\n",
			cmdChain,
			output.Tokens, 
			output.Metadata.Duration,
			lineCount,
			formatBytes(output.Metadata.Bytes))
	} else {
		fmt.Printf("ctx → %s | ✗ Failed (exit %d) | %d tokens | %d ms | %d lines | %s\n",
			cmdChain,
			output.Metadata.ExitCode, 
			output.Tokens,
			output.Metadata.Duration,
			lineCount,
			formatBytes(output.Metadata.Bytes))
	}
	
	// Separator
	fmt.Println("────────────────────────────────────────")
	
	// Output
	if output.Output != "" {
		fmt.Println(output.Output)
	}
	
	// Error details if present
	if output.Metadata.Error != "" {
		fmt.Printf("\nError: %s\n", output.Metadata.Error)
	}
	
	return nil
}

// ExecuteStreamCommand executes a command in streaming mode
func (ce *CommandExecutor) ExecuteStreamCommand(ctx context.Context, args []string) error {
	command := strings.Join(args, " ")
	
	// Apply timeout from configuration if set
	if ce.appCtx.Config.DefaultTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, ce.appCtx.Config.DefaultTimeout)
		defer cancel()
	}
	
	// Create streaming callback function
	lineCb := func(line string, streamType string) {
		event := models.StreamEvent{
			Type: streamType,
			Line: line,
		}
		data, err := json.Marshal(event)
		if err == nil {
			fmt.Println(string(data))
		}
	}
	
	// Get tokenizer for limit checking
	tok, _ := ce.appCtx.GetTokenizer()
	
	// Get limits from config
	var maxBytes, maxLines, maxTokens int64
	if ce.appCtx.Config.Limits.MaxOutputBytes != nil {
		maxBytes = *ce.appCtx.Config.Limits.MaxOutputBytes
	}
	if ce.appCtx.Config.Limits.MaxLines != nil {
		maxLines = *ce.appCtx.Config.Limits.MaxLines
	}
	if ce.appCtx.Config.Limits.MaxTokens != nil {
		maxTokens = *ce.appCtx.Config.Limits.MaxTokens
	}
	
	result, err := executor.ExecuteCommandStreaming(ctx, command, lineCb, tok, maxBytes, maxLines, maxTokens)
	if err != nil {
		// Output error as a stream event
		errorEvent := models.StreamEvent{
			Type: "stderr",
			Line: err.Error(),
		}
		data, _ := json.Marshal(errorEvent)
		fmt.Println(string(data))
		
		// Output final result envelope with error
		output := models.NewOutput(command, result.Output, result.ExitCode, result.Duration)
		output.Metadata.Error = err.Error()
		output.Metadata.Success = false
		
		// Set failure reason based on error type
		switch err {
		case executor.ErrLineLimitExceeded:
			output.Metadata.FailureReason = "line_limit_exceeded"
		case executor.ErrOutputLimitExceeded:
			output.Metadata.FailureReason = "output_limit_exceeded"
		case executor.ErrTokenLimitExceeded:
			output.Metadata.FailureReason = "token_limit_exceeded"
		}
		
		finalEvent := models.StreamEvent{
			Type:     "result",
			Envelope: output,
		}
		data, _ = json.Marshal(finalEvent)
		fmt.Println(string(data))
		
		return &ExitError{Code: ExitCodeWrappedCmdError}
	}
	
	// Enrich the final output (note: data.output will be empty for streaming)
	output, err := ce.enricher.EnrichOutput(ctx, result)
	if err != nil {
		return fmt.Errorf("failed to enrich output: %w", err)
	}
	
	// Post-execution token check for streaming
	if ce.appCtx.Config.MaxTokens > 0 && int64(output.Tokens) > ce.appCtx.Config.MaxTokens {
		limitErr := &TokenLimitExceededError{
			Limit:  ce.appCtx.Config.MaxTokens,
			Actual: output.Tokens,
		}
		output.Metadata.Error = limitErr.Error()
		output.Metadata.Success = false
		output.Metadata.FailureReason = "token_limit_exceeded"
	}
	
	// Clear the output since it was already streamed
	output.Output = ""
	
	// Output final result envelope
	finalEvent := models.StreamEvent{
		Type:     "result",
		Envelope: output,
	}
	data, err := json.Marshal(finalEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal final result: %w", err)
	}
	fmt.Println(string(data))
	
	// Return appropriate exit code
	if !output.Metadata.Success {
		// If the command itself failed, prioritize its exit code
		if result.ExitCode != 0 {
			return &ExitError{Code: ExitCodeWrappedCmdError}
		}
		// Otherwise, it was our token limit check that failed it
		return &ExitError{Code: ExitCodeWrappedCmdError}
	}
	
	if result.ExitCode != 0 {
		return &ExitError{Code: ExitCodeWrappedCmdError}
	}
	
	return nil
}