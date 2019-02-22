// +build !windows

package nodos

import (
	"io"
	"os"
)

func coInitializeEx(res uintptr, opt uintptr) {}

func coUninitialize() {}

func isEscapeSequenceAvailable() bool {
	return true
}

func getConsole() io.Writer {
	return os.Stdout
}
