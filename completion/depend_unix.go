// +build !windows

package completion

import (
	"os"
	"path/filepath"
)

func join(dir, name string) string {
	return filepath.Join(dir, name)
}

func isExecutable(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return (stat.Mode().Perm() & 0555) != 0
}
