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

func volumeName(drive *uint16) (label string, fsname string, err error) {
	var _label [256]uint16
	var serial, length, flags uint32
	var _fsname [256]uint16

	err = windows.GetVolumeInformation(drive, &_label[0], uint32(len(_label)), &serial, &length, &flags, &_fsname[0], uint32(len(_fsname)))
	if err != nil {
		return
	}
	label = windows.UTF16ToString(_label[:])
	fsname = windows.UTF16ToString(_fsname[:])
	return
}

func VolumeName(drive string) (label string, fsname string, err error) {
	_drive, err := windows.UTF16PtrFromString(drive)
	if err != nil {
		return "", "", err
	}
	return volumeName(_drive)
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
