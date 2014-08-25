package mbcs

/*
#include <windows.h>
*/
import "C"
import "syscall"
import "unsafe"

func UtoA(utf8 string) []byte {
	utf16, _ := syscall.UTF16FromString(utf8)
	size := C.WideCharToMultiByte(C.CP_THREAD_ACP, 0,
		(*C.WCHAR)(&utf16[0]), C.int(len(utf16)), nil, 0, nil, nil)
	mbcs := make([]byte, size)
	C.WideCharToMultiByte(C.CP_THREAD_ACP, 0,
		(*C.WCHAR)(&utf16[0]), C.int(len(utf16)),
		(*C.CHAR)(unsafe.Pointer(&mbcs[0])), size, nil, nil)
	return mbcs
}

func AtoU(mbcs []byte) string {
	size := C.MultiByteToWideChar(C.CP_THREAD_ACP, 0,
		(*C.CHAR)(unsafe.Pointer(&mbcs[0])), C.int(len(mbcs)), nil, 0)
	utf16 := make([]uint16, size)
	C.MultiByteToWideChar(C.CP_THREAD_ACP, 0,
		(*C.CHAR)(unsafe.Pointer(&mbcs[0])), C.int(len(mbcs)),
		(*C.WCHAR)(unsafe.Pointer(&utf16[0])), size)
	return syscall.UTF16ToString(utf16)
}
