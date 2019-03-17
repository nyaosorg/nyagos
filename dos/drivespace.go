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

type Drive struct {
	Letter rune
	Type   uint32
}

func GetDrives() ([]*Drive, error) {
	bits, err := windows.GetLogicalDrives()
	if err != nil {
		return nil, err
	}
	result := []*Drive{}
	for d := 'A'; d <= 'Z'; d++ {
		if (bits & 1) != 0 {
			rootPathName := []uint16{uint16(d), ':', '\\', 0}
			type1 := &Drive{Letter: d, Type: windows.GetDriveType(&rootPathName[0])}
			result = append(result, type1)
		}
		bits >>= 1
	}
	return result, nil
}
