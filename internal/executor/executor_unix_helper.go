//go:build !windows
// +build !windows

package executor

import (
	"os/exec"
)

// associateProcessWithJobObject is a no-op on Unix systems
// since they use process groups instead of Job Objects
func associateProcessWithJobObject(cmd *exec.Cmd) error {
	// Unix systems don't need this - process groups are set before starting
	return nil
}
