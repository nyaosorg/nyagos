package dos

import (
	"os"
	"path/filepath"
)

func truepath_(path string) string {
	parent := filepath.Dir(path)
	if parent != "" && parent != "." && parent != path {
		path = filepath.Join(truepath_(parent), filepath.Base(path))
	}
	if newpath, err := os.Readlink(path); err == nil {
		return newpath
	} else {
		return path
	}
}

func TruePath(path string) string {
	if newpath, err := filepath.Abs(path); err == nil {
		return truepath_(newpath)
	} else {
		return truepath_(newpath)
	}
}
