package nodos

import (
	"golang.org/x/sys/windows"
	"unsafe"

	"github.com/zetamatta/nyagos/dos"
)

var kernel32 = windows.NewLazyDLL("kernel32")

var procCopyFileW = kernel32.NewProc("CopyFileW")

// Copy calls Win32's CopyFile API.
func copyFile(src, dst string, isFailIfExists bool) error {
	_src, err := windows.UTF16PtrFromString(src)
	if err != nil {
		return err
	}
	_dst, err := windows.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}
	var _isFailIfExists uintptr
	if isFailIfExists {
		_isFailIfExists = 1
	} else {
		_isFailIfExists = 0
	}
	rc, _, err := procCopyFileW.Call(
		uintptr(unsafe.Pointer(_src)),
		uintptr(unsafe.Pointer(_dst)),
		_isFailIfExists)
	if rc == 0 {
		return err
	}
	return nil
}

// Move calls Win32's MoveFileEx API.
func moveFile(src, dst string) error {
	_src, err := windows.UTF16PtrFromString(src)
	if err != nil {
		return err
	}
	_dst, err := windows.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}
	return windows.MoveFileEx(
		_src,
		_dst,
		windows.MOVEFILE_REPLACE_EXISTING|
			windows.MOVEFILE_COPY_ALLOWED|
			windows.MOVEFILE_WRITE_THROUGH)
}

func readShortcut(path string) (string, string, error) {
	return dos.ReadShortcut(path)
}
