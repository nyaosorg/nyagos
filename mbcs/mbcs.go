package mbcs

/*
#include <windows.h>
*/
import "C"
import "syscall"
import "unsafe"
import "fmt"

func lastError(name string) error {
	switch C.GetLastError() {
	case C.ERROR_INSUFFICIENT_BUFFER:
		return fmt.Errorf("%s: ERROR_INSUFFICIENT_BUFFER", name)
	case C.ERROR_INVALID_FLAGS:
		return fmt.Errorf("%s: ERROR_INVALID_FLAGS", name)
	case C.ERROR_INVALID_PARAMETER:
		return fmt.Errorf("%s: ERROR_INVALID_PARAMETER", name)
	case C.ERROR_NO_UNICODE_TRANSLATION:
		return fmt.Errorf("%s: ERROR_NO_UNICODE_TRANSLATION", name)
	default:
		return fmt.Errorf("%s: Unknown error", name)
	}
}

func UtoA(utf8 string) ([]byte, error) {
	utf16, err := syscall.UTF16FromString(utf8)
	if err != nil {
		return nil, err
	}
	size := C.WideCharToMultiByte(C.CP_THREAD_ACP, 0,
		(*C.WCHAR)(&utf16[0]), C.int(len(utf16)), nil, 0, nil, nil)
	if size <= 0 {
		return nil, lastError("WideCharToMultiByte")
	}
	mbcs := make([]byte, size)
	rc := C.WideCharToMultiByte(C.CP_THREAD_ACP, 0,
		(*C.WCHAR)(&utf16[0]), C.int(len(utf16)),
		(*C.CHAR)(unsafe.Pointer(&mbcs[0])), size, nil, nil)
	if rc == 0 {
		return nil, lastError("WideCharToMultiByte")
	}
	return mbcs, nil
}

func AtoU(mbcs []byte) (string, error) {
	if mbcs == nil || len(mbcs) <= 0 {
		return "", nil
	}
	size := C.MultiByteToWideChar(C.CP_THREAD_ACP, 0,
		(*C.CHAR)(unsafe.Pointer(&mbcs[0])), C.int(len(mbcs)), nil, 0)
	if size <= 0 {
		return "", lastError("MultiByteToWideChar")
	}
	utf16 := make([]uint16, size)
	rc := C.MultiByteToWideChar(C.CP_THREAD_ACP, 0,
		(*C.CHAR)(unsafe.Pointer(&mbcs[0])), C.int(len(mbcs)),
		(*C.WCHAR)(&utf16[0]), size)
	if rc == 0 {
		return "", lastError("MultiByteToWideChar")
	}
	return syscall.UTF16ToString(utf16), nil
}
