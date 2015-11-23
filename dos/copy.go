package dos

import (
	"syscall"
	"unsafe"
)

var copyFile = kernel32.NewProc("CopyFileW")
var moveFileEx = kernel32.NewProc("MoveFileExW")

func Copy(src string, dst string, isFailIfExists bool) error {
	cSrc, cSrcErr := syscall.UTF16PtrFromString(src)
	if cSrcErr != nil {
		return cSrcErr
	}
	cDst, cDstErr := syscall.UTF16PtrFromString(dst)
	if cDstErr != nil {
		return cDstErr
	}
	var isFailIfExists_ uintptr
	if isFailIfExists {
		isFailIfExists_ = 1
	} else {
		isFailIfExists_ = 0
	}
	rc, _, err := copyFile.Call(
		uintptr(unsafe.Pointer(cSrc)),
		uintptr(unsafe.Pointer(cDst)),
		isFailIfExists_)
	if rc != 0 {
		return nil
	} else {
		return err
	}
}

func Move(src, dst string) error {
	cSrc, cSrcErr := syscall.UTF16PtrFromString(src)
	if cSrcErr != nil {
		return cSrcErr
	}
	cDst, cDstErr := syscall.UTF16PtrFromString(dst)
	if cDstErr != nil {
		return cDstErr
	}
	rc, _, err := moveFileEx.Call(
		uintptr(unsafe.Pointer(cSrc)),
		uintptr(unsafe.Pointer(cDst)),
		MOVEFILE_REPLACE_EXISTING|MOVEFILE_COPY_ALLOWED|MOVEFILE_WRITE_THROUGH)
	if rc != 0 {
		return nil
	} else {
		return err
	}
}
