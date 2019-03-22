package dos

import (
	"fmt"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var mpr = syscall.NewLazyDLL("mpr")
var wNetGetConnectionW = mpr.NewProc("WNetGetConnectionW")

func WNetGetConnection(localName string) (string, error) {
	localNamePtr, localNameErr := syscall.UTF16PtrFromString(localName)
	if localNameErr != nil {
		return "", localNameErr
	}
	var buffer [1024]uint16
	size := uintptr(len(buffer))

	rc, _, err := wNetGetConnectionW.Call(
		uintptr(unsafe.Pointer(localNamePtr)),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)))

	if uint32(rc) != 0 {
		return "", err
	}
	return syscall.UTF16ToString(buffer[:]), nil
}

type NetDrive struct {
	Letter rune
	Remote string
}

func GetNetDrives() ([]*NetDrive, error) {
	drives, err := GetDrives()
	if err != nil {
		return nil, err
	}
	result := []*NetDrive{}
	for _, d := range drives {
		if d.Type == windows.DRIVE_REMOTE {
			path, err := WNetGetConnection(fmt.Sprintf("%c:", d.Letter))
			if err == nil {
				node := &NetDrive{Letter: d.Letter, Remote: path}
				result = append(result, node)
			}
		}
	}
	return result, nil
}

// https://msdn.microsoft.com/ja-jp/library/cc447030.aspx
// http://eternalwindows.jp/security/share/share06.html
