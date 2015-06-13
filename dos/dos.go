package dos

import (
	"fmt"
	"regexp"
	"strings"
	"syscall"
	"unicode"
	"unsafe"
)

var msvcrt = syscall.NewLazyDLL("msvcrt")
var _chdrive = msvcrt.NewProc("_chdrive")
var _wchdir = msvcrt.NewProc("_wchdir")

func chdrive_(n rune) uintptr {
	rc, _, _ := _chdrive.Call(uintptr(n & 0x1F))
	return rc
}

func getFirst(s string) (rune, error) {
	reader := strings.NewReader(s)
	drive, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return unicode.ToUpper(drive), nil
}

// Change drive without changing the working directory there.
func Chdrive(drive string) error {
	driveLetter, driveErr := getFirst(drive)
	if driveErr != nil {
		return driveErr
	}
	chdrive_(driveLetter)
	return nil
}

var rxPath = regexp.MustCompile("^([a-zA-Z]):(.*)$")

// Change the current working directory
// without changeing the working directory
// in the last drive.
func Chdir(folder_ string) error {
	folder := folder_
	if m := rxPath.FindStringSubmatch(folder_); m != nil {
		status := chdrive_(rune(m[1][0]))
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
