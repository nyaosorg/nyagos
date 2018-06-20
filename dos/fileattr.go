package dos

import (
	"syscall"
)

// GetFileAttributes calls Win32-API's GetFileAttributes.
func GetFileAttributes(path string) (uint32, error) {
	cpath, cpathErr := syscall.UTF16PtrFromString(path)
	if cpathErr != nil {
		return 0, cpathErr
	}
	return syscall.GetFileAttributes(cpath)
}

// SetFileAttributes calls Win32-API's SetFileAttributes
func SetFileAttributes(path string, attr uint32) error {
	cpath, cpathErr := syscall.UTF16PtrFromString(path)
	if cpathErr != nil {
		return cpathErr
	}
	return syscall.SetFileAttributes(cpath, attr)
}
