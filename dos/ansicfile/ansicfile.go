package ansicfile

import (
	"syscall"
	"unsafe"
)

var msvcrt = syscall.NewLazyDLL("msvcrt")
var wfopen = msvcrt.NewProc("_wfopen")
var fclose = msvcrt.NewProc("fclose")
var fputc = msvcrt.NewProc("fputc")

func Open(path string, mode string) (uintptr, error) {
	path_ptr, path_err := syscall.UTF16PtrFromString(path)
	if path_err != nil {
		return 0, path_err
	}
	mode_ptr, mode_err := syscall.UTF16PtrFromString(mode)
	if mode_err != nil {
		return 0, mode_err
	}
	rc, _, err := wfopen.Call(uintptr(unsafe.Pointer(path_ptr)),
		uintptr(unsafe.Pointer(mode_ptr)))
	if rc == 0 {
		return 0, err
	} else {
		return rc, nil
	}
}

// Close has to be not a method but a simple function
// because it uses as callback function from DLLs (ex lua53.dll).

func Close(fp uintptr) {
	fclose.Call(fp)
}

func Putc(c byte, fp uintptr) {
	fputc.Call(uintptr(c), fp)
}
