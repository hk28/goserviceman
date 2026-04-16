//go:build windows

package process

import (
	"os/exec"
	"strconv"
)

func killProc(entry *procEntry) error {
	return exec.Command("taskkill", "/F", "/PID", strconv.Itoa(entry.PID)).Run()
}
