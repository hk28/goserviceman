//go:build windows

package browser

import "os/exec"

// Open opens url in the default system browser.
func Open(url string) error {
	return exec.Command("cmd", "/c", "start", "", url).Start()
}
