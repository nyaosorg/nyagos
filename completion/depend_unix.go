//go:build !windows
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
	mode := stat.Mode()
	return mode.IsRegular() && (mode.Perm()&0111) != 0
}
