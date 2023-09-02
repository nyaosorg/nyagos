//go:build !windows
// +build !windows

package nodos

import (
	"github.com/nyaosorg/go-windows-findfile"
)

func lookPathSkip(f *findfile.FileInfo) bool {
	return f.IsDir()
}
