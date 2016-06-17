package dos

import (
	"syscall"
)

func GetFileAttributes(path string) (uint32, error) {
	cpath, cpathErr := syscall.UTF16PtrFromString(path)
	if cpathErr != nil {
		return 0, cpathErr
	}
	return syscall.GetFileAttributes(cpath)
}

func SetFileAttributes(path string, attr uint32) error {
	cpath, cpathErr := syscall.UTF16PtrFromString(path)
	if cpathErr != nil {
		return cpathErr
	}
	return syscall.SetFileAttributes(cpath, attr)
}
