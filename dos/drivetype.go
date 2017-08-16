package dos

import (
	"syscall"
	"unsafe"
)

var procGetDriveType = kernel32.NewProc("GetDriveTypeW")

func GetDriveType(rootPathName string) (uintptr, error) {
	path, err := syscall.UTF16PtrFromString(rootPathName)
	if err != nil {
		return 0, err
	}
	rc, _, err := procGetDriveType.Call(uintptr(unsafe.Pointer(path)))
	if rc == DRIVE_UNKNOWN || rc == DRIVE_NO_ROOT_DIR {
		return 0, err
	}
	return rc, nil
}
