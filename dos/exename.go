package dos

import (
	"os"
	"syscall"
	"unsafe"
)

var getModuleFileName = kernel32.NewProc("GetModuleFileNameW")

func GetModuleFileName() (string, error) {
	var path16 [syscall.MAX_PATH]uint16
	result, _, err := getModuleFileName.Call(0, uintptr(unsafe.Pointer(&path16[0])), uintptr(len(path16)))
	if result == 0 {
		return os.Args[0], err
	}
	return syscall.UTF16ToString(path16[:]), nil
}
