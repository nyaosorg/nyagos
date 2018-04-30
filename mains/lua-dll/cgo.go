package lua

import (
	"unsafe"
)

func CGoBytes(p, length uintptr) []byte {
	if length <= 0 || p == 0 {
		return []byte{}
	}
	buffer := make([]byte, length)
	copyMemory(uintptr(unsafe.Pointer(&buffer[0])), p, length)
	return buffer
}

func CGoStringN(p, length uintptr) string {
	if length <= 0 || p == 0 {
		return ""
	}
	return string(CGoBytes(p, length))
}

func CGoStringZ(p uintptr) string {
	return CGoStringN(p, strLen(p))
}
