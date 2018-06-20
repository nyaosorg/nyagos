package commands

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

func correct(path string) (string, string, error) {
	dirname, fname := filepath.Split(filepath.Clean(path))
	fi, err := ioutil.ReadDir(dirname)
	if err != nil {
		return dirname, fname, err
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

// CorrectCase corrects `path`'s case.
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
