package dos

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

var kernel32 = windows.NewLazySystemDLL("kernel32.dll")
var procGetDiskFreeSpaceExW = kernel32.NewProc("GetDiskFreeSpaceExW")

// GetDiskFreeSpace retunrs disk information.
//   rootPathName - string like "C:"
func GetDiskFreeSpace(rootPathName string) (free uint64, total uint64, totalFree uint64, err error) {
	_rootPathName, err := windows.UTF16PtrFromString(rootPathName)
	if err != nil {
		return 0, 0, 0, err
	}
	rc, _, err := procGetDiskFreeSpaceExW.Call(
		uintptr(unsafe.Pointer(_rootPathName)),
		uintptr(unsafe.Pointer(&free)),
		uintptr(unsafe.Pointer(&total)),
		uintptr(unsafe.Pointer(&totalFree)))
	if rc == 0 {
		return 0, 0, 0, err
	}
	return free, total, totalFree, nil
}
