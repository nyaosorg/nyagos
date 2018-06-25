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

// Chdrive changes drive without changing the working directory there.
func Chdrive(drive string) error {
	for _, c := range drive {
		chDriveSub(unicode.ToUpper(c))
		return nil
	}
	return errors.New("Chdrive: driveletter not found")
}

var rxPath = regexp.MustCompile("^([a-zA-Z]):(.*)$")

// Chdir changes the current working directory
// without changeing the working directory
// in the last drive.
func Chdir(_folder string) error {
	folder := _folder
	if m := rxPath.FindStringSubmatch(_folder); m != nil {
		status := chDriveSub(rune(m[1][0]))
		if status != 0 {
			return fmt.Errorf("%s: no such directory", _folder)
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
			err = fmt.Errorf("%s: no such directory", _folder)
		}
	}
	return err
}
