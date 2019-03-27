package dos

import (
	"errors"
	"unsafe"

	"golang.org/x/sys/windows"
)

var mpr = windows.NewLazyDLL("mpr")
var procWNetGetConnectionW = mpr.NewProc("WNetGetConnectionW")

func _WNetGetConnection(drive uint16) (string, error) {
	localName := []uint16{drive, ':', 0}
	var buffer [1024]uint16
	size := uintptr(len(buffer))

	rc, _, err := procWNetGetConnectionW.Call(
		uintptr(unsafe.Pointer(&localName[0])),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)))

	if uint32(rc) != 0 {
		return "", err
	}
	return windows.UTF16ToString(buffer[:]), nil
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
			path, err := _WNetGetConnection(uint16(d.Letter))
			if err == nil {
				node := &NetDrive{Letter: d.Letter, Remote: path}
				result = append(result, node)
			}
		}
	}
	return result, nil
}

func FindVacantDrive() (uint, error) {
	bits, err := windows.GetLogicalDrives()
	if err != nil {
		return 0, err
	}
	for d := uint('Z'); d >= 'A'; d-- {
		if (bits & (1 << (d - 'A'))) == 0 {
			return d, nil
		}
	}
	return 0, errors.New("vacant drive is not found")
}

// https://msdn.microsoft.com/ja-jp/library/cc447030.aspx
// http://eternalwindows.jp/security/share/share06.html
