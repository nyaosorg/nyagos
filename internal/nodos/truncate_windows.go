package nodos

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

// Truncate is same as os.RemoveAll but report files to remove.
func truncate(folder string, whenError func(string, error) bool, out io.Writer) error {
	attr, err := GetFileAttributes(folder)
	if err != nil {
		return fmt.Errorf("%s: %w", folder, err)
	}
	if (attr & REPARSE_POINT) == 0 {
		// Only not junction, delete files under folder.
		files, err := os.ReadDir(folder)
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
				SetFileAttributes(fullpath, windows.FILE_ATTRIBUTE_NORMAL)
				err = windows.Unlink(fullpath)
			}
			if err != nil {
				if whenError != nil && !whenError(fullpath, err) {
					return fmt.Errorf("%s: %s", fullpath, err.Error())
				}
			}
		}
	}
	if (attr & windows.FILE_ATTRIBUTE_READONLY) != 0 {
		SetFileAttributes(folder, attr&^windows.FILE_ATTRIBUTE_READONLY)
	}
	if err := windows.Rmdir(folder); err != nil {
		return fmt.Errorf("%s: %w", folder, err)
	}
	return nil
}
