package lua

import "syscall"

var msvcrt = syscall.NewLazyDLL("msvcrt")
var procMemcpy = msvcrt.NewProc("memcpy")
var procStrlen = msvcrt.NewProc("strlen")

func copyMemory(dst uintptr, src uintptr, length uintptr) {
	procMemcpy.Call(dst, src, length)
}

func strLen(p uintptr) uintptr {
	rc, _, _ := procStrlen.Call(p)
	return rc
}
