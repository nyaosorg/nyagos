package dos

import "fmt"
import "regexp"
import "strings"
import "unicode"
import "syscall"
import "path/filepath"

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

func Chdir(folder_ string) error {
	folder := folder_
	if m := rxPath.FindStringSubmatch(folder_); m != nil {
		status := C._chdrive(C.int(m[1][0] & 0x1F))
		if status != 0 {
			return fmt.Errorf("%s: no such directory", folder_)
		}
		folder = m[2]
		if len(folder) <= 0 {
			return nil
		}
	}
	utf16, err := syscall.UTF16FromString(folder)
	if err == nil {
		status := C._wchdir((*C.wchar_t)(&utf16[0]))
		if status != 0 {
			err = fmt.Errorf("%s: no such directory", folder_)
		}
	}
	return err
}

var rxDriveOnly = regexp.MustCompile("^[a-zA-Z]:$")
var rxRoot = regexp.MustCompile("^([a-zA-Z]:)?[\\/]")

func Join(paths ...string) string {
	start := 0
	for i, path := range paths {
		if rxDriveOnly.MatchString(path) {
			paths[i] = path + "."
		} else if rxRoot.MatchString(path) {
			start = i
		}
	}
	if start > 0 {
		paths = paths[start:]
	}
	return filepath.Join(paths...)
}

var rxCouldGlobPattern = regexp.MustCompile("^[A-Za-z]:[^\\/]")

func Glob(pattern string) (matches []string, err error) {
	result, err := filepath.Glob(pattern)
	if len(result) > 0 {
		return result, err
	}
	if rxCouldGlobPattern.MatchString(pattern) {
		pattern = fmt.Sprintf("%s.\\%s", pattern[:2], pattern[2:])
		result, err = filepath.Glob(pattern)
	}
	return result, err
}
