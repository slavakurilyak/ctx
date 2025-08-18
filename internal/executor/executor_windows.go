//go:build windows
// +build windows

package executor

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	
	procCreateJobObjectW         = kernel32.NewProc("CreateJobObjectW")
	procAssignProcessToJobObject = kernel32.NewProc("AssignProcessToJobObject")
	procSetInformationJobObject  = kernel32.NewProc("SetInformationJobObject")
	procTerminateJobObject       = kernel32.NewProc("TerminateJobObject")
	procCloseHandle              = kernel32.NewProc("CloseHandle")
)

const (
	// Job Object limit flags
	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE = 0x00002000
	
	// Job Object information classes
	JobObjectExtendedLimitInformation = 9
)

// JOBOBJECT_BASIC_LIMIT_INFORMATION structure
type JOBOBJECT_BASIC_LIMIT_INFORMATION struct {
	PerProcessUserTimeLimit int64
	PerJobUserTimeLimit     int64
	LimitFlags              uint32
	MinimumWorkingSetSize   uintptr
	MaximumWorkingSetSize   uintptr
	ActiveProcessLimit      uint32
	Affinity                uintptr
	PriorityClass           uint32
	SchedulingClass         uint32
}

// JOBOBJECT_EXTENDED_LIMIT_INFORMATION structure
type JOBOBJECT_EXTENDED_LIMIT_INFORMATION struct {
	BasicLimitInformation JOBOBJECT_BASIC_LIMIT_INFORMATION
	IoInfo                IO_COUNTERS
	ProcessMemoryLimit    uintptr
	JobMemoryLimit        uintptr
	PeakProcessMemoryUsed uintptr
	PeakJobMemoryUsed     uintptr
}

// IO_COUNTERS structure
type IO_COUNTERS struct {
	ReadOperationCount  uint64
	WriteOperationCount uint64
	OtherOperationCount uint64
	ReadTransferCount   uint64
	WriteTransferCount  uint64
	OtherTransferCount  uint64
}

// jobObjectMap stores Job Object handles associated with each command
// This allows concurrent command execution with separate job objects
var (
	jobObjectMap = make(map[*exec.Cmd]syscall.Handle)
	jobObjectMux sync.Mutex
)

// setupProcessGroup configures the command to use a Windows Job Object
func setupProcessGroup(cmd *exec.Cmd) {
	if cmd.SysProcAttr == nil {
		cmd.SysProcAttr = &syscall.SysProcAttr{}
	}
	
	// Create a new Job Object for this command
	job, err := createJobObject()
	if err != nil {
		// Fall back to default behavior if Job Object creation fails
		// This might happen on older Windows versions or restricted environments
		return
	}
	
	// Store the job handle associated with this command
	jobObjectMux.Lock()
	jobObjectMap[cmd] = job
	jobObjectMux.Unlock()
	
	// Set CREATE_NEW_PROCESS_GROUP flag to allow Ctrl+C handling
	cmd.SysProcAttr.CreationFlags = syscall.CREATE_NEW_PROCESS_GROUP
}

// createJobObject creates a new Windows Job Object
func createJobObject() (syscall.Handle, error) {
	// Create an unnamed job object
	ret, _, err := procCreateJobObjectW.Call(
		0, // lpJobAttributes (NULL for default security)
		0, // lpName (NULL for unnamed object)
	)
	
	if ret == 0 {
		return 0, fmt.Errorf("CreateJobObject failed: %v", err)
	}
	
	job := syscall.Handle(ret)
	
	// Configure the job object to terminate all processes when the job handle is closed
	limitInfo := JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	
	ret, _, err = procSetInformationJobObject.Call(
		uintptr(job),
		uintptr(JobObjectExtendedLimitInformation),
		uintptr(unsafe.Pointer(&limitInfo)),
		unsafe.Sizeof(limitInfo),
	)
	
	if ret == 0 {
		procCloseHandle.Call(uintptr(job))
		return 0, fmt.Errorf("SetInformationJobObject failed: %v", err)
	}
	
	return job, nil
}

// associateProcessWithJob associates a process with a Job Object
func associateProcessWithJob(job syscall.Handle, processHandle syscall.Handle) error {
	ret, _, err := procAssignProcessToJobObject.Call(
		uintptr(job),
		uintptr(processHandle),
	)
	
	if ret == 0 {
		return fmt.Errorf("AssignProcessToJobObject failed: %v", err)
	}
	
	return nil
}

// killProcessGroup terminates all processes in the Job Object
func killProcessGroup(cmd *exec.Cmd) error {
	jobObjectMux.Lock()
	job, exists := jobObjectMap[cmd]
	jobObjectMux.Unlock()
	
	if exists && job != 0 {
		// Terminate all processes in the job object
		ret, _, err := procTerminateJobObject.Call(
			uintptr(job),
			1, // Exit code
		)
		
		if ret == 0 {
			// If termination fails, try to kill the process directly as fallback
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
			// Don't return error, just continue to cleanup
		}
		
		// Close the job object handle
		procCloseHandle.Call(uintptr(job))
		
		// Remove from map
		jobObjectMux.Lock()
		delete(jobObjectMap, cmd)
		jobObjectMux.Unlock()
		
		return nil
	}
	
	// Fallback: kill the process directly if no job object
	if cmd.Process != nil {
		return cmd.Process.Kill()
	}
	
	return nil
}

// configureTermination sets up termination behavior for Windows
func configureTermination(cmd *exec.Cmd) {
	// Override the default Cancel function to use Job Object termination
	cmd.Cancel = func() error {
		return killProcessGroup(cmd)
	}
	
	// Set WaitDelay to give processes time to cleanup
	// Default to 3 seconds, but allow configuration via environment variable
	waitDelay := 3 * time.Second
	if v := os.Getenv("CTX_WAIT_DELAY"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			waitDelay = d
		}
	}
	cmd.WaitDelay = waitDelay
}

// AssociateWithJobObject associates the command's process with its job object
// This should be called after cmd.Start() but before any child processes are created
func AssociateWithJobObject(cmd *exec.Cmd) error {
	if cmd.Process == nil {
		return fmt.Errorf("process not started")
	}
	
	jobObjectMux.Lock()
	job, exists := jobObjectMap[cmd]
	jobObjectMux.Unlock()
	
	if exists && job != 0 {
		// Get the process handle
		// On Windows, cmd.Process.Pid is the process ID, we need to get the handle
		processHandle, err := syscall.OpenProcess(syscall.PROCESS_ALL_ACCESS, false, uint32(cmd.Process.Pid))
		if err != nil {
			return fmt.Errorf("OpenProcess failed: %v", err)
		}
		defer syscall.CloseHandle(processHandle)
		
		return associateProcessWithJob(job, processHandle)
	}
	
	return nil
}