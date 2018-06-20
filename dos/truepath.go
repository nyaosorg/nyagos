package dos

import (
	"os"
	"path/filepath"
)

func truePathSub(path string) string {
	parent := filepath.Dir(path)
	if parent != "" && parent != "." && parent != path {
		path = filepath.Join(truePathSub(parent), filepath.Base(path))
	}
	if linkTo, err := os.Readlink(path); err == nil {
		path = linkTo
	}
	return path
}

// TruePath returns the casesensitive filename
func TruePath(path string) string {
	if newpath, err := filepath.Abs(path); err == nil {
		path = newpath
	}
	return truePathSub(path)
}
