// +build !windows

package completion

import (
	"os"
)

func isExecutable(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}
	return (stat.Mode().Perm() & 0555) != 0
}
