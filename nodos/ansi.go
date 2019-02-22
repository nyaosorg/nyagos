package nodos

import (
	"io"
)

func IsEscapeSequenceAvailable() bool {
	return isEscapeSequenceAvailable()
}

func GetConsole() io.Writer {
	return getConsole()
}
