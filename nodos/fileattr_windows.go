package nodos

import (
	"golang.org/x/sys/windows"
)

const (
	reparsePoint = windows.FILE_ATTRIBUTE_REPARSE_POINT
)

// GetFileAttributes calls Win32-API's GetFileAttributes.
func getFileAttributes(path string) (uint32, error) {
	cpath, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}
	attr, err := windows.GetFileAttributes(cpath)
	return attr, err
}

// SetFileAttributes calls Win32-API's SetFileAttributes
func setFileAttributes(path string, attr uint32) error {
	cpath, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return windows.SetFileAttributes(cpath, attr)
}
