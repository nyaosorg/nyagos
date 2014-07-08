package dos

import "regexp"
import "strings"
import "unicode"
import "syscall"

//#include <direct.h>
import "C"

func GetFirst(s string) (rune, error) {
	reader := strings.NewReader(s)
	drive, _, err := reader.ReadRune()
	if err != nil {
		return 0, err
	}
	return unicode.ToUpper(drive), nil
}

func Chdrive(drive string) error {
	driveLetter, driveErr := GetFirst(drive)
	if driveErr != nil {
		return driveErr
	}
	C._chdrive(C.int(driveLetter) & 0x1F)
	return nil
}

var rxPath = regexp.MustCompile("^([a-zA-Z]):(.*)$")

func Chdir(folder string) error {
	if m := rxPath.FindStringSubmatch(folder); m != nil {
		C._chdrive(C.int(m[1][0] & 0x1F))
		folder = m[2]
	}
	utf16, err := syscall.UTF16FromString(folder)
	if err == nil {
		C._wchdir((*C.wchar_t)(&utf16[0]))
	}
	return err
}
