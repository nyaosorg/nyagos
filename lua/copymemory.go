package lua

import "syscall"

var kernel32 = syscall.NewLazyDLL("kernel32")
var procCopyMemory = kernel32.NewProc("RtlCopyMemory")

func copyMemory(dst uintptr, src uintptr, length uintptr) {
	procCopyMemory.Call(dst, src, length)
}

var msvcrt = syscall.NewLazyDLL("msvcrt")
var procStrlen = msvcrt.NewProc("strlen")

func strLen(p uintptr) uintptr {
	rc, _, _ := procStrlen.Call(p)
	return rc
}
