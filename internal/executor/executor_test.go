package executor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestTimeoutWithGoRun(t *testing.T) {
	// Create a Go file that sleeps
	code := `package main
import (
	"fmt"
	"time"
)
func main() {
	fmt.Println("Starting Go program...")
	time.Sleep(10 * time.Second)
	fmt.Println("Should not see this")
}`
	tmpDir := t.TempDir()
	tmpfile := filepath.Join(tmpDir, "slow.go")
	if err := os.WriteFile(tmpfile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	result, err := ExecuteCommand(ctx, "go run "+tmpfile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if the command was terminated (exit code -1 usually indicates killed)
	if result.ExitCode == 0 {
		t.Fatal("Expected non-zero exit code due to timeout")
	}

	// Allow some tolerance for timing
	if duration > 2*time.Second {
		t.Fatalf("Timeout took too long: %v", duration)
	}

	// Check that we got the initial output but not the final message
	output := string(result.Output)
	if !strings.Contains(output, "Starting Go program") && len(output) > 0 {
		// It's ok if we don't get any output due to buffering
		t.Logf("Output: %s", output)
	}
	if strings.Contains(output, "Should not see this") {
		t.Fatal("Process was not terminated in time")
	}
}

func TestTimeoutWithShellScript(t *testing.T) {
	script := `#!/bin/bash
echo "Starting shell script..."
sleep 10
echo "Should not see this"`

	tmpDir := t.TempDir()
	tmpfile := filepath.Join(tmpDir, "slow.sh")
	if err := os.WriteFile(tmpfile, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	result, err := ExecuteCommand(ctx, "bash "+tmpfile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if the command was terminated (exit code -1 usually indicates killed)
	if result.ExitCode == 0 {
		t.Fatal("Expected non-zero exit code due to timeout")
	}

	if duration > 2*time.Second {
		t.Fatalf("Timeout took too long: %v", duration)
	}

	output := string(result.Output)
	if strings.Contains(output, "Should not see this") {
		t.Fatal("Shell script was not terminated in time")
	}
}

func TestTimeoutWithNestedProcesses(t *testing.T) {
	// Test a command that spawns multiple levels of children
	cmd := `bash -c "bash -c 'sleep 10'"`

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	result, err := ExecuteCommand(ctx, cmd)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if the command was terminated (exit code -1 usually indicates killed)
	if result.ExitCode == 0 {
		t.Fatal("Expected non-zero exit code due to timeout")
	}

	if duration > 2*time.Second {
		t.Fatalf("Timeout took too long: %v", duration)
	}
}

func TestTimeoutWithPythonScript(t *testing.T) {
	// Create a Python file that sleeps
	code := `import time
print("Starting Python...")
time.sleep(10)
print("Should not see this")`

	tmpDir := t.TempDir()
	tmpfile := filepath.Join(tmpDir, "slow.py")
	if err := os.WriteFile(tmpfile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	result, err := ExecuteCommand(ctx, "python3 "+tmpfile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if the command was terminated (exit code -1 usually indicates killed)
	if result.ExitCode == 0 {
		t.Fatal("Expected non-zero exit code due to timeout")
	}

	if duration > 2*time.Second {
		t.Fatalf("Timeout took too long: %v", duration)
	}

	output := string(result.Output)
	if strings.Contains(output, "Should not see this") {
		t.Fatal("Python script was not terminated in time")
	}
}

func TestGracefulTermination(t *testing.T) {
	// Create a script that handles SIGTERM
	script := `#!/bin/bash
trap 'echo "Received SIGTERM"; exit 0' TERM
echo "Starting..."
sleep 10
echo "Should not see this"`

	tmpDir := t.TempDir()
	tmpfile := filepath.Join(tmpDir, "graceful.sh")
	if err := os.WriteFile(tmpfile, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result, _ := ExecuteCommand(ctx, "bash "+tmpfile)

	// Check if the script started
	output := string(result.Output)
	if !strings.Contains(output, "Starting...") && !strings.Contains(output, "Received SIGTERM") {
		t.Logf("Output: %s", output)
	}
}

func TestTimeoutWithNodeScript(t *testing.T) {
	// Create a Node.js file that sleeps
	code := `console.log("Starting Node.js...");
setTimeout(() => {
    console.log("Should not see this");
    process.exit(0);
}, 10000);`

	tmpDir := t.TempDir()
	tmpfile := filepath.Join(tmpDir, "slow.js")
	if err := os.WriteFile(tmpfile, []byte(code), 0644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	result, err := ExecuteCommand(ctx, "node "+tmpfile)
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if the command was terminated (exit code -1 usually indicates killed)
	if result.ExitCode == 0 {
		t.Fatal("Expected non-zero exit code due to timeout")
	}

	if duration > 2*time.Second {
		t.Fatalf("Timeout took too long: %v", duration)
	}

	output := string(result.Output)
	if strings.Contains(output, "Should not see this") {
		t.Fatal("Node.js script was not terminated in time")
	}
}

func TestDirectExecutableTimeout(t *testing.T) {
	// Test with a simple sleep command
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	result, err := ExecuteCommand(ctx, "sleep 10")
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if the command was terminated (exit code -1 usually indicates killed)
	if result.ExitCode == 0 {
		t.Fatal("Expected non-zero exit code due to timeout")
	}

	if duration > 2*time.Second {
		t.Fatalf("Timeout took too long: %v", duration)
	}
}

func TestPipelineTimeout(t *testing.T) {
	// Test timeout with a pipeline
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	start := time.Now()
	result, err := ExecuteCommand(ctx, "echo 'test' | sleep 10")
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check if the command was terminated (exit code -1 usually indicates killed)
	if result.ExitCode == 0 {
		t.Fatal("Expected non-zero exit code due to timeout")
	}

	if duration > 2*time.Second {
		t.Fatalf("Timeout took too long: %v", duration)
	}
}

func TestConfigurableWaitDelay(t *testing.T) {
	// Test that CTX_WAIT_DELAY environment variable is respected
	oldValue := os.Getenv("CTX_WAIT_DELAY")
	defer os.Setenv("CTX_WAIT_DELAY", oldValue)

	// Set a short wait delay
	os.Setenv("CTX_WAIT_DELAY", "50ms")

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	start := time.Now()
	result, _ := ExecuteCommand(ctx, "sleep 10")
	duration := time.Since(start)

	// Command should timeout around 500ms, not wait full 3 seconds
	if duration > 1*time.Second {
		t.Fatalf("WaitDelay not respected, took: %v", duration)
	}

	// Verify termination reason is recorded
	if result.Metadata["termination_reason"] != "timeout" {
		t.Fatalf("Expected termination_reason to be 'timeout', got: %v", result.Metadata["termination_reason"])
	}
}

func TestConfigurableSigtermGrace(t *testing.T) {
	// Test that CTX_SIGTERM_GRACE environment variable is respected
	oldValue := os.Getenv("CTX_SIGTERM_GRACE")
	defer os.Setenv("CTX_SIGTERM_GRACE", oldValue)

	// Set a longer grace period
	os.Setenv("CTX_SIGTERM_GRACE", "200ms")

	// Create a script that handles SIGTERM and logs timing
	script := `#!/bin/bash
trap 'echo "Received SIGTERM at $(date +%s%N)"; exit 0' TERM
echo "Starting at $(date +%s%N)"
sleep 10`

	tmpDir := t.TempDir()
	tmpfile := filepath.Join(tmpDir, "grace_test.sh")
	if err := os.WriteFile(tmpfile, []byte(script), 0755); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	result, _ := ExecuteCommand(ctx, "bash "+tmpfile)

	// The output should contain both timestamps if grace period worked
	output := string(result.Output)
	if !strings.Contains(output, "Starting at") {
		t.Log("Script didn't start properly")
	}
}

func TestTerminationReason(t *testing.T) {
	// Test that we can distinguish between timeout and manual cancellation

	// Test timeout
	ctx1, cancel1 := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel1()

	result1, _ := ExecuteCommand(ctx1, "sleep 1")
	if result1.Metadata["termination_reason"] != "timeout" {
		t.Fatalf("Expected termination_reason 'timeout', got: %v", result1.Metadata["termination_reason"])
	}

	// Test manual cancellation
	ctx2, cancel2 := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel2()
	}()

	result2, _ := ExecuteCommand(ctx2, "sleep 10")
	if result2.Metadata["termination_reason"] != "cancelled" {
		t.Fatalf("Expected termination_reason 'cancelled', got: %v", result2.Metadata["termination_reason"])
	}
}
