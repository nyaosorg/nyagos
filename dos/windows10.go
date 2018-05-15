package dos

import (
	"unsafe"
)

// var kernel32 = syscall.NewLazyDLL("kernel32")

// const STD_INPUT_HANDLE = uintptr(1) + ^uintptr(10)
const STD_OUTPUT_HANDLE = uintptr(1) + ^uintptr(11)
const STD_ERROR_HANDLE = uintptr(1) + ^uintptr(12)
const ENABLE_VIRTUAL_TERMINAL_PROCESSING uintptr = 0x0004

var procGetStdHandle = kernel32.NewProc("GetStdHandle")
var procGetConsoleMode = kernel32.NewProc("GetConsoleMode")
var procSetConsoleMode = kernel32.NewProc("SetConsoleMode")

func enableVirtualTerminalProcessing(handle uintptr) (func(), error) {
	var mode uintptr
	console, _, _ := procGetStdHandle.Call(handle)
	rc, _, err := procGetConsoleMode.Call(console, uintptr(unsafe.Pointer(&mode)))
	if rc == 0 {
		return func() {}, err
	}
	deferFunc := func() { procSetConsoleMode.Call(console, mode) }

	rc, _, err = procSetConsoleMode.Call(console, mode|ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	if rc == 0 {
		return deferFunc, err
	}
	return deferFunc, nil
}

func EnableStdoutVirtualTerminalProcessing() (func(), error) {
	return enableVirtualTerminalProcessing(STD_OUTPUT_HANDLE)
}

func EnableStderrVirtualTerminalProcessing() (func(), error) {
	return enableVirtualTerminalProcessing(STD_ERROR_HANDLE)
}
