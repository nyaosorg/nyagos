package dos

import (
	"fmt"
	"os"
	"path/filepath"
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

func CorrectPathCase(path string) (string, error) {
	if len(path) <= 3 {
		return strings.ToUpper(path), nil
	}
	dirname, fname, err := correct(path)
	if err != nil {
		return path, err
	}
	if len(dirname) > 3 {
		// NOT root directory.
		dirname, err = CorrectPathCase(dirname)
	} else {
		// root directory.
		dirname = strings.ToUpper(dirname)
	}
	return filepath.Join(dirname, fname), nil
}

func Getwd() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	wd, _ = CorrectPathCase(wd)
	return wd, nil
}
