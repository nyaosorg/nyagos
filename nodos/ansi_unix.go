// +build !windows

package nodos

import (
	"io"
	"os"
)

func isEscapeSequenceAvailable() bool {
	return true
}

func getConsole() io.Writer {
	return os.Stdout
}
