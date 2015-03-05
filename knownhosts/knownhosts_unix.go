// Build known_hosts path on unix type machines.

// +build !windows

package knownhosts

import (
	"os"
	"path/filepath"
)

// path returns the path to the OpenSSH known_hosts file.
// This function is platform specific.
func path() string {
	home := os.Getenv("HOME")
	path := filepath.Join(home, ".ssh", "known_hosts")
	return path
}
