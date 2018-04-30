package dos

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"syscall"
)

func Truncate(folder string, whenError func(string, error) bool, out io.Writer) error {
	attr, err := GetFileAttributes(folder)
	if err != nil {
		return fmt.Errorf("%s: %s", folder, err)
	}
	if (attr & FILE_ATTRIBUTE_REPARSE_POINT) == 0 {
		// Only not junction, delete files under folder.
		files, err := ioutil.ReadDir(folder)
		if err != nil {
			return err
		}
		for _, f := range files {
			if f.Name() == "." || f.Name() == ".." {
				continue
			}
			fullpath := filepath.Join(folder, f.Name())
			var err error
			if f.IsDir() {
				fmt.Fprintf(out, "%s\\\n", fullpath)
				err = Truncate(fullpath, whenError, out)
			} else {
				fmt.Fprintln(out, fullpath)
				SetFileAttributes(fullpath, FILE_ATTRIBUTE_NORMAL)
				err = syscall.Unlink(fullpath)
			}
			if err != nil {
				if whenError != nil && !whenError(fullpath, err) {
					return fmt.Errorf("%s: %s", fullpath, err.Error())
				}
			}
		}
	}
	if err := syscall.Rmdir(folder); err != nil {
		return fmt.Errorf("%s: %s", folder, err.Error())
	}
	return nil
}
