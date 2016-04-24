package dos

import (
	"path/filepath"
	"syscall"
	"unsafe"
)

var mpr = syscall.NewLazyDLL("mpr")
var wNetGetConnectionW = mpr.NewProc("WNetGetConnectionW")

func WNetGetConnection(localName string) (string, error) {
	localNamePtr, localNameErr := syscall.UTF16PtrFromString(localName)
	if localNameErr != nil {
		return "", localNameErr
	}
	var buffer [1024]uint16
	var size uintptr = uintptr(len(buffer))

	rc, _, err := wNetGetConnectionW.Call(
		uintptr(unsafe.Pointer(localNamePtr)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)))

	if uint32(rc) != 0 {
		return "", err
	}
	return syscall.UTF16ToString(buffer[:]), nil
}

func NetDriveToUNC(path string) string {
	if path[1] == ':' {
		// print("'", path[:2], "'\n")
		path_, err := WNetGetConnection(path[:2])
		if err == nil {
			return filepath.Join(path_, path[2:])
		}
		// print(err.Error(), "\n")
	}
	return path
}
