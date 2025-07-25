//go:build !windows

package main

import (
	"os"
)

func sendInterruptSignal(pid int) error {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return proc.Signal(os.Interrupt)
}
