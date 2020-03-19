package dos

import (
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/sys/windows"

	"github.com/zetamatta/nyagos/nodos"
)

// Truncate is same as os.RemoveAll but report files to remove.
func Truncate(folder string, whenError func(string, error) bool, out io.Writer) error {
	attr, err := nodos.GetFileAttributes(folder)
	if err != nil {
		return fmt.Errorf("%s: %s", folder, err)
	}
	if (attr & nodos.REPARSE_POINT) == 0 {
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
				nodos.SetFileAttributes(fullpath, windows.FILE_ATTRIBUTE_NORMAL)
				err = windows.Unlink(fullpath)
			}
			if err != nil {
				if whenError != nil && !whenError(fullpath, err) {
					return fmt.Errorf("%s: %s", fullpath, err.Error())
				}
			}
		}
	}
	if err := windows.Rmdir(folder); err != nil {
		return fmt.Errorf("%s: %s", folder, err.Error())
	}
	return nil
}
