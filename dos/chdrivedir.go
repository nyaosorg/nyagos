package dos

import (
	"errors"
	"fmt"
	"regexp"
	"syscall"
	"unicode"
	"unsafe"
)

var msvcrt = syscall.NewLazyDLL("msvcrt")
var _chdrive = msvcrt.NewProc("_chdrive")
var _wchdir = msvcrt.NewProc("_wchdir")

func chDriveSub(n rune) uintptr {
	rc, _, _ := _chdrive.Call(uintptr(n & 0x1F))
	return rc
}

// Change drive without changing the working directory there.
func Chdrive(drive string) error {
	for _, c := range drive {
		chDriveSub(unicode.ToUpper(c))
		return nil
	}
	return errors.New("Chdrive: driveletter not found")
}

var rxPath = regexp.MustCompile("^([a-zA-Z]):(.*)$")

// Change the current working directory
// without changeing the working directory
// in the last drive.
func Chdir(folder_ string) error {
	folder := folder_
	if m := rxPath.FindStringSubmatch(folder_); m != nil {
		status := chDriveSub(rune(m[1][0]))
		if status != 0 {
			return fmt.Errorf("%s: no such directory", folder_)
		}
		folder = m[2]
		if len(folder) <= 0 {
			return nil
		}
	}
	utf16, err := syscall.UTF16PtrFromString(folder)
	if err == nil {
		status, _, _ := _wchdir.Call(uintptr(unsafe.Pointer(utf16)))
		if status != 0 {
			err = fmt.Errorf("%s: no such directory", folder_)
		}
	}
	return err
}
