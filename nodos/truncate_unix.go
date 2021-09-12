//go:build !windows
// +build !windows

package nodos

import (
	"io"
	"os"
)

func truncate(folder string, _ func(string, error) bool, _ io.Writer) error {
	return os.RemoveAll(folder)
}
