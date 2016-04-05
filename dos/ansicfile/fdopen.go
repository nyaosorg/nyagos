package ansicfile

import (
	"syscall"
	"unsafe"
)

var wfdopen = msvcrt.NewProc("_wfdopen")

func FdOpen(handle uintptr, mode string) (FilePtr, error) {
	mode_ptr, mode_err := syscall.UTF16PtrFromString(mode)
	if mode_err != nil {
		return 0, mode_err
	}
	rc, _, err := wfdopen.Call(handle, uintptr(unsafe.Pointer(mode_ptr)))
	if rc == 0 {
		return 0, err
	} else {
		return FilePtr(rc), nil
	}
}
