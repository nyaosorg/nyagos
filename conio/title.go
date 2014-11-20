package conio

import (
	"syscall"
	"unsafe"
)

var setConsoleTitle = kernel32.NewProc("SetConsoleTitleW")

func SetTitle(title string) {
	ctitle, err := syscall.UTF16FromString(title)
	if ctitle != nil && err == nil {
		setConsoleTitle.Call(uintptr(unsafe.Pointer(&ctitle[0])))
	}
}
