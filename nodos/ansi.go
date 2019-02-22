package nodos

import (
	"io"
)

func CoInitializeEx(res uintptr, opt uintptr) {
	coInitializeEx(res, opt)
}

func CoUninitialize() {
	coUninitialize()
}

func IsEscapeSequenceAvailable() bool {
	return isEscapeSequenceAvailable()
}

func GetConsole() io.Writer {
	return getConsole()
}
