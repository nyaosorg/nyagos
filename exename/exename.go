package exename

import "syscall"
import "unsafe"

var kernel32 = syscall.NewLazyDLL("kernel32")
var procGetModuleFileName = kernel32.NewProc("GetModuleFileNameW")

func Query() string {
	var path16 [syscall.MAX_PATH]uint16
	procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&path16[0])), uintptr(len(path16)))
	return syscall.UTF16ToString(path16[:])
}
