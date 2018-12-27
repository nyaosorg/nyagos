package dos

import (
	"fmt"
	"os"

	"golang.org/x/sys/windows"
)

func IsEscapeSequenceAvailable() bool {
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
