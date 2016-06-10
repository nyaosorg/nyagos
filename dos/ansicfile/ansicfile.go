package ansicfile

import (
	"errors"
	"syscall"
	"unsafe"
)

var msvcrt = syscall.NewLazyDLL("msvcrt")
var wfopen = msvcrt.NewProc("_wfopen")
var fclose = msvcrt.NewProc("fclose")
var fwrite = msvcrt.NewProc("fwrite")

type FilePtr uintptr

func Open(path string, mode string) (FilePtr, error) {
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
		return FilePtr(rc), nil
	}
}

func (fp FilePtr) Close() {
	fclose.Call(uintptr(fp))
}

func (fp FilePtr) Write(p []byte) (int, error) {
	rc, _, err := fwrite.Call(uintptr(unsafe.Pointer(&p[0])),
		1, uintptr(len(p)), uintptr(fp))
	n := int(rc)
	if n == len(p) {
		return n, nil
	} else if err != nil {
		return n, err
	} else {
		return n, errors.New("ansicfile.FilePtr.Write error")
	}
}
