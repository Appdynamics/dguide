//go:build !windows
// +build !windows

package roller

import (
	"golang.org/x/crypto/ssh/terminal"
)

// IsTerminal checks if the given file descriptor is a terminal.
func IsTerminal(fd int) bool {
	return terminal.IsTerminal(fd)
}
