package exename

import "os"
import "strings"
import "syscall"
import "unsafe"

var kernel32 = syscall.NewLazyDLL("kernel32")
var procGetModuleFileName = kernel32.NewProc("GetModuleFileNameW")

func Query() (string, error) {
	var path16 [syscall.MAX_PATH]uint16
	result, _, err := procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&path16[0])), uintptr(len(path16)))
	if result == 0 {
		return os.Args[0], err
	}
	return syscall.UTF16ToString(path16[:]), nil
}

var Suffixes = map[string]bool{}

func init() {
	pathExt := os.Getenv("PATHEXT")
	for _, ext := range strings.Split(pathExt, ";") {
		Suffixes[strings.ToLower(ext)] = true
	}
}
