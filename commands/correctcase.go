package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func correct(path string) (string, string, error) {
	dirname, fname := filepath.Split(filepath.Clean(path))
	fd, fdErr := os.Open(dirname)
	if fdErr != nil {
		return dirname, fname, fdErr
	}
	defer fd.Close()
	fi, fiErr := fd.Readdir(-1)
	if fiErr != nil {
		return dirname, fname, fiErr
	}
	for _, fi1 := range fi {
		if strings.EqualFold(fi1.Name(), fname) {
			return dirname, fi1.Name(), nil
		}
	}
	return dirname, fname, fmt.Errorf("%s: not found.", path)
}

var rxRoot1 = regexp.MustCompile(`^[a-zA-Z]:\\?$`)
var rxRoot2 = regexp.MustCompile(`^\\\\\w+\\\w+$`)

// correct path's case.
func CorrectCase(path string) (string, error) {
	if rxRoot1.MatchString(path) {
		return strings.ToUpper(path), nil
	}
	if rxRoot2.MatchString(path) {
		return path, nil
	}
	dirname, fname, err := correct(path)
	if err != nil {
		return path, err
	}
	if len(dirname) > 0 {
		// NOT root directory.
		dirname, _ = CorrectCase(dirname)
	}
	return filepath.Join(dirname, fname), nil
}
