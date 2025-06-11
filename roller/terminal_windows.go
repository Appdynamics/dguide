//go:build windows
// +build windows

package roller

import (
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// IsTerminal checks if the given file descriptor is a terminal.
func IsTerminal(fd syscall.Handle) bool {
	return terminal.IsTerminal(int(fd))
}
