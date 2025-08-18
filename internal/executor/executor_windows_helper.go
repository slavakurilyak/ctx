//go:build windows
// +build windows

package executor

import (
	"os/exec"
)

// associateProcessWithJobObject is a Windows-specific helper function
// that associates a started process with its Job Object
func associateProcessWithJobObject(cmd *exec.Cmd) error {
	return AssociateWithJobObject(cmd)
}
