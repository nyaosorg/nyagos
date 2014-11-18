package dos

import "fmt"
import "io"
import "os"
import "syscall"

func Truncate(folder string, out io.Writer) error {
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
		if f.IsDir() {
			fmt.Fprintf(out, "%s\\\n", fullpath)
			if err := Truncate(fullpath, out); err != nil {
				return fmt.Errorf("%s: %s", fullpath, err.Error())
			}
		} else {
			fmt.Fprintln(out, fullpath)
			syscall.Unlink(fullpath)
		}
	}
	if err := syscall.Rmdir(folder); err != nil {
		return fmt.Errorf("%s: %s", folder, err.Error())
	}
	return nil
}
