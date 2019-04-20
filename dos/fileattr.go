package dos

import (
	"golang.org/x/sys/windows"
)

// GetFileAttributes calls Win32-API's GetFileAttributes.
func GetFileAttributes(path string) (uint32, error) {
	cpath, cpathErr := windows.UTF16PtrFromString(path)
	if cpathErr != nil {
		return 0, cpathErr
	}
	return windows.GetFileAttributes(cpath)
}

// SetFileAttributes calls Win32-API's SetFileAttributes
func SetFileAttributes(path string, attr uint32) error {
	cpath, cpathErr := windows.UTF16PtrFromString(path)
	if cpathErr != nil {
		return cpathErr
	}
	return windows.SetFileAttributes(cpath, attr)
}
