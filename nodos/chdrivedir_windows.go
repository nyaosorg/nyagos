package nodos

import (
	"golang.org/x/sys/windows"

	dos "github.com/zetamatta/go-windows-netresource"
)

func chdriveRetry(c rune) bool {
	drive := []uint16{uint16(c), ':', 0}
	t := windows.GetDriveType(&drive[0])
	if t != windows.DRIVE_REMOTE {
		return false
	}
	uncpath, err := dos.WNetGetConnectionUTF16a(uint16(c))
	if err != nil {
		return false
	}
	if _, err := dos.NetUse(uncpath, string([]byte{byte(c), ':', 0})); err != nil {
		return false
	}
	return true
}
