package nodos

import (
	"github.com/mattn/go-colorable"
	"io"
)

func CoInitializeEx(res uintptr, opt uintptr) {
	coInitializeEx(res, opt)
}

func CoUninitialize() {
	coUninitialize()
}

func GetConsole() io.Writer {
	return colorable.NewColorableStdout()
}
