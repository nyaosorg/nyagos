package dos

import (
	"os"
	"syscall"

	"github.com/zetamatta/go-findfile"
)

func GetFileAttributesFromFileInfo(status os.FileInfo) uint32 {
	if it, ok := status.Sys().(*syscall.Win32FileAttributeData); ok {
		return it.FileAttributes
	} else if it, ok := status.(*findfile.FileInfo); ok {
		return it.FileAttributes
	} else {
		panic("Can not get fileatttribute")
	}
}

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
