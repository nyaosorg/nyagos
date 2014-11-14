package dos

import "syscall"
import "unsafe"

var copyFile = kernel32.NewProc("CopyFileW")
var moveFile = kernel32.NewProc("MoveFileW")

func Copy(src string, dst string, override bool) error {
	cSrc, cSrcErr := syscall.UTF16FromString(src)
	if cSrcErr != nil {
		return cSrcErr
	}
	cDst, cDstErr := syscall.UTF16FromString(dst)
	if cDstErr != nil {
		return cDstErr
	}
	var override_ uintptr
	if override {
		override_ = 1
	} else {
		override_ = 0
	}
	rc, _, err := copyFile.Call(
		uintptr(unsafe.Pointer(&cSrc[0])),
		uintptr(unsafe.Pointer(&cDst[0])),
		override_)
	if rc != 0 {
		return nil
	} else {
		return err
	}
}

func Move(src, dst string) error {
	cSrc, cSrcErr := syscall.UTF16FromString(src)
	if cSrcErr != nil {
		return cSrcErr
	}
	cDst, cDstErr := syscall.UTF16FromString(dst)
	if cDstErr != nil {
		return cDstErr
	}
	rc, _, err := moveFile.Call(
		uintptr(unsafe.Pointer(&cSrc[0])),
		uintptr(unsafe.Pointer(&cDst[0])))
	if rc != 0 {
		return nil
	} else {
		return err
	}
}
