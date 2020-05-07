package nodos

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"
)

// chDriveRune changes drive and returns the previous working directory.
func chDriveRune(n rune) (string, error) {
	lastDir, lastDirErr := os.Getwd()

	newDir := os.Getenv(fmt.Sprintf("=%c:", n))
	if newDir == "" {
		newDir = fmt.Sprintf("%c:%c", n, os.PathSeparator)
	}
	err := os.Chdir(newDir)
	if err != nil {
		return lastDir, err
	}
	if lastDirErr == nil {
		os.Setenv("="+filepath.VolumeName(lastDir), lastDir)
	}
	return lastDir, err
}

// Chdrive changes drive without changing the working directory there.
// And returns the previous working directory.
func Chdrive(drive string) (string, error) {
	c, _ := utf8.DecodeRuneInString(drive)
	return chDriveRune(c)
}

var rxPath = regexp.MustCompile("^([a-zA-Z]):(.*)$")

// Chdir changes the current working directory
// without changeing the working directory
// in the last drive.
func Chdir(folder string) (err error) {
	lastDir := ""
	if m := rxPath.FindStringSubmatch(folder); m != nil {
		folder = m[2]
		lastDir, err = chDriveRune(rune(m[1][0]))
		if err != nil {
			return err
		}
		if len(folder) <= 0 { // Change drive only.
			return nil
		}
	}
	absFolder, absErr := filepath.Abs(folder)
	err = os.Chdir(folder)
	if err != nil {
		if lastDir != "" {
			os.Chdir(lastDir)
		}
		return err
	}
	if absErr == nil {
		folder = absFolder
	}
	os.Setenv("="+filepath.VolumeName(folder), folder)
	return nil
}
