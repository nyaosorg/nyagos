package nodos

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

var shell32 = windows.NewLazySystemDLL("shell32.dll")
var procSHOpenWithDialog = shell32.NewProc("SHOpenWithDialog")

type _OpenAsInfo struct {
	FileName *uint16
	Class    *uint16
	Flag     uintptr
}

func shOpenWithDialog(filename, class string) (err error) {
	var info _OpenAsInfo

	info.FileName, err = windows.UTF16PtrFromString(filename)
	if err != nil {
		return
	}
	if class != "" {
		info.Class, err = windows.UTF16PtrFromString(class)
		if err != nil {
			return
		}
	}
	info.Flag = 0x4
	var rc uintptr
	rc, _, err = procSHOpenWithDialog.Call(0, uintptr(unsafe.Pointer(&info)))
	if rc != S_OK {
		return
	}
	return nil
}
