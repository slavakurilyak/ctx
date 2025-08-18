package executor

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/slavakurilyak/ctx/internal/tokenizer"
)

// Define error types for limit exceeded conditions
var (
	ErrLineLimitExceeded   = errors.New("line limit exceeded")
	ErrOutputLimitExceeded = errors.New("output limit exceeded")
	ErrTokenLimitExceeded  = errors.New("token limit exceeded")
)

type ExecutionResult struct {
	Output   []byte
	ExitCode int
	Duration time.Duration
	Command  string
	Metadata map[string]interface{}
}

func ExecuteCommand(ctx context.Context, command string) (*ExecutionResult, error) {
	start := time.Now()

	// Detect if we need shell wrapping for complex commands
	needsShell := strings.ContainsAny(command, "|<>&;`$")

	var cmd *exec.Cmd
	if needsShell {
		// Use shell for complex commands
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		cmd = exec.CommandContext(ctx, shell, "-c", command)
	} else {
		// Direct execution for simple commands
		parts := splitCommand(command)
		if len(parts) == 0 {
			return nil, fmt.Errorf("empty command")
		}
		cmd = exec.CommandContext(ctx, parts[0], parts[1:]...)
	}

	// Setup process group for proper cleanup (platform-specific)
	setupProcessGroup(cmd)

	// Configure termination behavior (platform-specific)
	configureTermination(cmd)

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Inherit environment variables
	cmd.Env = os.Environ()

	// Set working directory to current directory
	if wd, err := os.Getwd(); err == nil {
		cmd.Dir = wd
	}

	// Start the command (non-blocking)
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	// On Windows, associate the process with the Job Object
	// On Unix, this is a no-op since process groups are set before starting
	if err := associateProcessWithJobObject(cmd); err != nil {
		// Log the error but continue - process will still run, just without job object protection
		// In production, you might want to handle this differently
		_ = err
	}

	// Wait for the command to complete
	err = cmd.Wait()

	// If context was cancelled, ensure process group is killed
	if ctx.Err() == context.DeadlineExceeded || ctx.Err() == context.Canceled {
		_ = killProcessGroup(cmd) // Ensure complete cleanup
	}

	duration := time.Since(start)

	// Determine exit code
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	// Combine stdout and stderr
	output := stdout.Bytes()
	if stderr.Len() > 0 {
		if len(output) > 0 {
			output = append(output, '\n')
		}
		output = append(output, stderr.Bytes()...)
	}

	result := &ExecutionResult{
		Output:   output,
		ExitCode: exitCode,
		Duration: duration,
		Command:  command,
		Metadata: make(map[string]interface{}),
	}

	// Add termination reason to metadata
	if ctx.Err() == context.DeadlineExceeded {
		result.Metadata["termination_reason"] = "timeout"
	} else if ctx.Err() == context.Canceled {
		result.Metadata["termination_reason"] = "cancelled"
	}

	return result, nil
}

func splitCommand(command string) []string {
	var parts []string
	var current []rune
	var inQuote rune
	var escaped bool

	for _, char := range command {
		if escaped {
			current = append(current, char)
			escaped = false
			continue
		}

		if char == '\\' {
			escaped = true
			continue
		}

		if inQuote != 0 {
			if char == inQuote {
				inQuote = 0
			} else {
				current = append(current, char)
			}
			continue
		}

		if char == '"' || char == '\'' {
			inQuote = char
			continue
		}

		if char == ' ' || char == '\t' {
			if len(current) > 0 {
				parts = append(parts, string(current))
				current = []rune{}
			}
			continue
		}

		current = append(current, char)
	}

	if len(current) > 0 {
		parts = append(parts, string(current))
	}

	return parts
}

// ExecuteCommandStreaming executes a command and streams its output line by line.
// It returns the final ExecutionResult after the command completes.
// The function now accepts limits for bytes, lines, and tokens to terminate early if exceeded.
func ExecuteCommandStreaming(
	ctx context.Context,
	command string,
	lineCb func(line string, streamType string),
	tok tokenizer.Tokenizer,
	maxBytes int64,
	maxLines int64,
	maxTokens int64,
) (*ExecutionResult, error) {
	start := time.Now()

	// Create a cancellable context for early termination on limit exceeded
	cmdCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Detect if we need shell wrapping for complex commands
	needsShell := strings.ContainsAny(command, "|<>&;`$")

	var cmd *exec.Cmd
	if needsShell {
		// Use shell for complex commands
		shell := os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
		cmd = exec.CommandContext(cmdCtx, shell, "-c", command)
	} else {
		// Direct execution for simple commands
		parts := splitCommand(command)
		if len(parts) == 0 {
			return nil, fmt.Errorf("empty command")
		}
		cmd = exec.CommandContext(cmdCtx, parts[0], parts[1:]...)
	}

	// Setup process group for proper cleanup (platform-specific)
	setupProcessGroup(cmd)

	// Configure termination behavior (platform-specific)
	configureTermination(cmd)

	// Set up pipes for streaming
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Inherit environment variables
	cmd.Env = os.Environ()

	// Set working directory to current directory
	if wd, err := os.Getwd(); err == nil {
		cmd.Dir = wd
	}

	var totalBytes int64
	var totalLines int64
	var totalTokens int64
	var wg sync.WaitGroup
	var stdoutBuf, stderrBuf bytes.Buffer
	var limitErr error
	var limitErrMux sync.Mutex

	wg.Add(2)
	go func() {
		if err := streamPipe(stdoutPipe, "stdout", &totalBytes, &totalLines, &totalTokens, &wg, lineCb, &stdoutBuf, tok, maxBytes, maxLines, maxTokens); err != nil {
			limitErrMux.Lock()
			limitErr = err
			limitErrMux.Unlock()
			cancel() // Cancel the command execution
		}
	}()
	go func() {
		if err := streamPipe(stderrPipe, "stderr", &totalBytes, &totalLines, &totalTokens, &wg, lineCb, &stderrBuf, tok, maxBytes, maxLines, maxTokens); err != nil {
			limitErrMux.Lock()
			if limitErr == nil {
				limitErr = err
			}
			limitErrMux.Unlock()
			cancel() // Cancel the command execution
		}
	}()

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// On Windows, associate the process with the Job Object
	// On Unix, this is a no-op since process groups are set before starting
	if err := associateProcessWithJobObject(cmd); err != nil {
		// Log the error but continue - process will still run, just without job object protection
		_ = err
	}

	// Wait for pipes to be fully read
	wg.Wait()

	// Wait for command to complete
	err = cmd.Wait()

	// If context was cancelled, ensure process group is killed
	if cmdCtx.Err() == context.DeadlineExceeded || cmdCtx.Err() == context.Canceled {
		_ = killProcessGroup(cmd) // Ensure complete cleanup
	}

	duration := time.Since(start)

	// Check if we had a limit error
	limitErrMux.Lock()
	finalLimitErr := limitErr
	limitErrMux.Unlock()

	// Determine exit code
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = 1
		}
	}

	// Combine stdout and stderr for the final result
	output := stdoutBuf.Bytes()
	if stderrBuf.Len() > 0 {
		if len(output) > 0 {
			output = append(output, '\n')
		}
		output = append(output, stderrBuf.Bytes()...)
	}

	result := &ExecutionResult{
		Output:   output,
		ExitCode: exitCode,
		Duration: duration,
		Command:  command,
		Metadata: make(map[string]interface{}),
	}

	// Add termination reason to metadata
	if cmdCtx.Err() == context.DeadlineExceeded {
		result.Metadata["termination_reason"] = "timeout"
	} else if cmdCtx.Err() == context.Canceled {
		result.Metadata["termination_reason"] = "cancelled"
	}

	// Return the limit error if one occurred
	if finalLimitErr != nil {
		return result, finalLimitErr
	}

	return result, nil
}

func streamPipe(
	pipe io.ReadCloser,
	streamType string,
	byteCounter *int64,
	lineCounter *int64,
	tokenCounter *int64,
	wg *sync.WaitGroup,
	cb func(string, string),
	buf *bytes.Buffer,
	tok tokenizer.Tokenizer,
	maxBytes int64,
	maxLines int64,
	maxTokens int64,
) error {
	defer wg.Done()
	defer pipe.Close()

	scanner := bufio.NewScanner(pipe)
	for scanner.Scan() {
		line := scanner.Text()

		// Check line limit
		if maxLines > 0 {
			newLineCount := atomic.AddInt64(lineCounter, 1)
			if newLineCount > maxLines {
				return ErrLineLimitExceeded
			}
		}

		// Check byte limit
		newByteCount := atomic.AddInt64(byteCounter, int64(len(line)+1)) // +1 for newline
		if maxBytes > 0 && newByteCount > maxBytes {
			return ErrOutputLimitExceeded
		}

		// Check token limit
		if tok != nil && maxTokens > 0 {
			count, err := tok.CountTokens(line)
			if err == nil { // Best-effort tokenization
				newTokenCount := atomic.AddInt64(tokenCounter, int64(count))
				if newTokenCount > maxTokens {
					return ErrTokenLimitExceeded
				}
			}
		}

		// Write to buffer for final result
		buf.WriteString(line + "\n")

		// Call the streaming callback
		cb(line, streamType)
	}

	return nil
}
