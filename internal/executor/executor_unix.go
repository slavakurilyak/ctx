//go:build !windows
// +build !windows

package executor

import (
	"os"
	"os/exec"
	"syscall"
	"time"
)

// setupProcessGroup configures the command to create a new process group
func setupProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	// Create new process group to ensure all child processes can be terminated
	cmd.SysProcAttr.Setpgid = true
}

// killProcessGroup terminates the entire process group
func killProcessGroup(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return nil
	}
	
	// Send SIGTERM to entire process group (negative PID targets the group)
	// This gives processes a chance to cleanup gracefully
	if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
		// If SIGTERM fails, immediately try SIGKILL
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL); err != nil {
			// Process might already be dead, which is fine
			return nil
		}
	}
	
	// Give processes a brief moment to handle SIGTERM
	// Default to 100ms, but allow configuration via environment variable
	gracePeriod := 100 * time.Millisecond
	if v := os.Getenv("CTX_SIGTERM_GRACE"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			gracePeriod = d
		}
	}
	time.Sleep(gracePeriod)
	
	// Force kill any remaining processes in the group
	// Ignore errors as processes might have already terminated
	_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	
	return nil
}

// configureTermination sets up proper cancellation for process groups
func configureTermination(cmd *exec.Cmd) {
	// When context is cancelled, kill the entire process group
	cmd.Cancel = func() error {
		return killProcessGroup(cmd)
	}
	// Give processes time to cleanup after SIGTERM before SIGKILL
	// Default to 3 seconds, but allow configuration via environment variable
	waitDelay := 3 * time.Second
	if v := os.Getenv("CTX_WAIT_DELAY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			waitDelay = d
		}
	}
	cmd.WaitDelay = waitDelay
}