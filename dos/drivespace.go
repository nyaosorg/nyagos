package dos

import (
	"syscall"
	"unsafe"
)

var procGetDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")

func GetDiskFreeSpace(rootPathName string) (free uint64, total uint64, totalFree uint64, err error) {
	path, err1 := syscall.UTF16PtrFromString(rootPathName)
	if err1 != nil {
		err = err1
		return
	}
	rc, _, err1 := procGetDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(path)),
		uintptr(unsafe.Pointer(&free)),
		uintptr(unsafe.Pointer(&total)),
		uintptr(unsafe.Pointer(&totalFree)))
	if rc == 0 {
		err = err1
	}
	return
}
