package dos

import (
	"syscall"
	"unsafe"
)

var getLogicalDrive = kernel32.NewProc("GetLogicalDrives")

func GetLogicalDrives() []rune {
	result := make([]rune, 0, 26)
	bits, _, _ := getLogicalDrive.Call()
	for i := 0; i < 26; i++ {
		if (bits & 1) != 0 {
			result = append(result, 'A'+rune(i))
		}
		bits >>= 1
	}
	return result
}

var getDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")

func GetDiskFreeSpaceEx(rootPath string) (int64, int64, int64, error) {
	var freeBytesAvailable [2]int64
	var totalNumberOfBytes [2]int64
	var totalNumberOfFreeBytes [2]int64

	rootPathW, rootPathErr := syscall.UTF16PtrFromString(rootPath)
	if rootPathErr != nil {
		return 0, 0, 0, rootPathErr
	}
	getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(rootPathW)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)))
	return freeBytesAvailable[0],
		totalNumberOfBytes[0],
		totalNumberOfFreeBytes[0],
		nil
}

/*
func Df() map[rune][2]int64{
	drives := GetLogicalDrives()
	for ,d := range drives {
		drive := fmt.Sprintf("%c:",d)
		GetDiskFreeSpace

	}
}
*/
