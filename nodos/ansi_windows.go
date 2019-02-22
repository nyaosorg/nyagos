package nodos

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/windows"

	"github.com/mattn/go-colorable"

	"github.com/zetamatta/nyagos/dos"
)

var ole32 = windows.NewLazyDLL("ole32")
var procCoInitializeEx = ole32.NewProc("CoInitializeEx")
var procCoUninitialize = ole32.NewProc("CoUninitialize")

func coInitializeEx(res uintptr, opt uintptr) {
	procCoInitializeEx.Call(res, opt)
}

func coUninitialize() {
	procCoUninitialize.Call()
}

func isEscapeSequenceAvailable() bool {
	var mode uint32
	err := windows.GetConsoleMode(windows.Stdout, &mode)
	if err != nil {
		return false
	}
	err = windows.SetConsoleMode(windows.Stdout, mode|0x4)
	if err != nil {
		return false
	}

	fmt.Print("\r\x1B[11G")
	os.Stdout.Sync()
	var csbi windows.ConsoleScreenBufferInfo

	err = windows.GetConsoleScreenBufferInfo(windows.Stdout, &csbi)
	result := (err == nil && csbi.CursorPosition.X == 10)
	fmt.Print("\r     \r")
	os.Stdout.Sync()
	windows.SetConsoleMode(windows.Stdout, mode)
	return result
}

var console io.Writer

var isEscapeSequenceAvailableFlag = false

func getConsole() io.Writer {
	if isEscapeSequenceAvailableFlag {
		dos.EnableStdoutVirtualTerminalProcessing()
		console = os.Stdout
	} else if console == nil {
		if isEscapeSequenceAvailable() {
			console = os.Stdout
			dos.EnableStdoutVirtualTerminalProcessing()
			isEscapeSequenceAvailableFlag = true
		} else {
			console = colorable.NewColorableStdout()
		}
	}
	return console
}
