package dos

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"
)

func chDriveRune(n rune) error {
	lastDir, lastDirErr := os.Getwd()

	newDir := os.Getenv(fmt.Sprintf("=%c:", n))
	if newDir == "" {
		newDir = fmt.Sprintf("%c:%c", n, os.PathSeparator)
	}
	err := os.Chdir(newDir)
	if err == nil && lastDirErr == nil {
		os.Setenv("="+filepath.VolumeName(lastDir), lastDir)
	}
	return err
}

// Chdrive changes drive without changing the working directory there.
func Chdrive(drive string) error {
	c, _ := utf8.DecodeRuneInString(drive)
	return chDriveRune(c)
}

var rxPath = regexp.MustCompile("^([a-zA-Z]):(.*)$")

// Chdir changes the current working directory
// without changeing the working directory
// in the last drive.
func Chdir(folder string) error {
	if m := rxPath.FindStringSubmatch(folder); m != nil {
		err := chDriveRune(rune(m[1][0]))
		if err != nil {
			return err
		}
		folder = m[2]
		if len(folder) <= 0 {
			return nil
		}
	}
	err := os.Chdir(folder)
	if err == nil {
		if absFolder, err := filepath.Abs(folder); err == nil {
			folder = absFolder
		}
		os.Setenv("="+filepath.VolumeName(folder), folder)
	}
	return err
}
