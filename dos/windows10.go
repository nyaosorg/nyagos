package dos

import (
	"golang.org/x/sys/windows"
)

const enableVirtualTerminalProcessing uint32 = 0x0004

func doEnableVirtualTerminalProcessing(console windows.Handle) (func(), error) {
	var mode uint32
	err := windows.GetConsoleMode(console, &mode)
	if err != nil {
		return func() {}, err
	}
	deferFunc := func() { windows.SetConsoleMode(console, mode) }
	err = windows.SetConsoleMode(console, mode|enableVirtualTerminalProcessing)
	return deferFunc, err
}

// EnableStdoutVirtualTerminalProcessing enables Windows10's native ESCAPE SEQUENCE support on STDOUT
func EnableStdoutVirtualTerminalProcessing() (func(), error) {
	return doEnableVirtualTerminalProcessing(windows.Stdout)
}

// EnableStderrVirtualTerminalProcessing enables Windows10's native ESCAPE SEQUENCE support on STDERR
func EnableStderrVirtualTerminalProcessing() (func(), error) {
	return doEnableVirtualTerminalProcessing(windows.Stderr)
}
