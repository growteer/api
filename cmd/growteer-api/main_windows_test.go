//go:build windows

package main

import (
	"golang.org/x/sys/windows"
)

func sendInterruptSignal(pid int) error {
	// NOTE (c.nicola): PROCESS_TERMINATE is more or less the equivalent to signal.SIGTERM.
	// Windows does not have interrupt signals (at least not in Go), so
	// we need to use the windows module and terminate the process.
	// This won't trigger the graceful shutdown routines, so keep that in mind.
	proc, err := windows.OpenProcess(windows.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return err
	}

	if err = windows.TerminateProcess(proc, 0); err != nil {
		return err
	}

	return windows.CloseHandle(proc)
}
