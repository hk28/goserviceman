//go:build !windows

package process

import (
	"os"
	"syscall"
)

func killProc(entry *procEntry) error {
	if entry.Cmd != nil && entry.Cmd.Process != nil {
		return entry.Cmd.Process.Signal(syscall.SIGTERM)
	}
	proc, err := os.FindProcess(entry.PID)
	if err != nil {
		return err
	}
	return proc.Signal(syscall.SIGTERM)
}
