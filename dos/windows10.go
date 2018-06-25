package dos

import (
	"unsafe"
)

// var kernel32 = syscall.NewLazyDLL("kernel32")

// const stdInputHandle = uintptr(1) + ^uintptr(10)
const stdOutputHandle = uintptr(1) + ^uintptr(11)
const stdErrorHandle = uintptr(1) + ^uintptr(12)
const enableVirtualTerminalProcessing uintptr = 0x0004

var procGetStdHandle = kernel32.NewProc("GetStdHandle")
var procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
var procSetConsoleMode = kernel32.NewProc("SetConsoleMode")

func doEnableVirtualTerminalProcessing(handle uintptr) (func(), error) {
	var mode uintptr
	console, _, _ := procGetStdHandle.Call(handle)
	rc, _, err := procGetConsoleMode.Call(console, uintptr(unsafe.Pointer(&mode)))
	if rc == 0 {
		return func() {}, err
	}
	deferFunc := func() { procSetConsoleMode.Call(console, mode) }

	rc, _, err = procSetConsoleMode.Call(console, mode|enableVirtualTerminalProcessing)
	if rc == 0 {
		return deferFunc, err
	}
	return deferFunc, nil
}

// EnableStdoutVirtualTerminalProcessing enables Windows10's native ESCAPE SEQUENCE support on STDOUT
func EnableStdoutVirtualTerminalProcessing() (func(), error) {
	return doEnableVirtualTerminalProcessing(stdOutputHandle)
}

// EnableStderrVirtualTerminalProcessing enables Windows10's native ESCAPE SEQUENCE support on STDERR
func EnableStderrVirtualTerminalProcessing() (func(), error) {
	return doEnableVirtualTerminalProcessing(stdErrorHandle)
}
