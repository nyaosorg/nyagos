package dos

import (
	"fmt"
	"io"
	"os"
	"syscall"
)

func Truncate(folder string, whenError func(string, error) bool, out io.Writer) error {
	attr, attrErr := GetFileAttributes(folder)
	if attrErr != nil {
		return fmt.Errorf("%s: %s", folder, attrErr)
	}
	if (attr & FILE_ATTRIBUTE_REPARSE_POINT) == 0 {
		// Only not junction, delete files under folder.
		fd, fdErr := os.Open(folder)
		if fdErr != nil {
			return fdErr
		}
		fi, fiErr := fd.Readdir(-1)
		fd.Close()
		if fiErr != nil {
			return fiErr
		}
		for _, f := range fi {
			if f.Name() == "." || f.Name() == ".." {
				continue
			}
			fullpath := Join(folder, f.Name())
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
