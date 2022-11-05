package nodos

import (
	"fmt"
	"io"
	"path/filepath"

	"golang.org/x/sys/windows"

	"github.com/nyaosorg/go-windows-findfile"
)

// Truncate is same as os.RemoveAll but report files to remove.
func truncate(folder string, whenError func(string, error) bool, out io.Writer) error {
	attr, err := GetFileAttributes(folder)
	if err != nil {
		return fmt.Errorf("%s: GetFileAttributes: %w", folder, err)
	}
	if (attr & REPARSE_POINT) == 0 {
		// Only not junction, delete files under folder.
		subDirectories := []string{}
		err1 := findfile.Walk(filepath.Join(folder, "*"), func(f *findfile.FileInfo) bool {
			if f.Name() == "." || f.Name() == ".." {
				return true
			}
			fullpath := filepath.Join(folder, f.Name())
			if f.IsDir() || f.IsReparsePoint() {
				subDirectories = append(subDirectories, fullpath)
			} else {
				fmt.Fprintln(out, fullpath)
				SetFileAttributes(fullpath, windows.FILE_ATTRIBUTE_NORMAL)
				if err = windows.Unlink(fullpath); err != nil {
					if whenError != nil && !whenError(fullpath, err) {
						err = fmt.Errorf("%s: windows.Unlink: %w", fullpath, err)
						return false
					}
				}
			}
			return true
		})
		if err1 != nil {
			return fmt.Errorf("%s: findfile.Walk: %w", folder, err1)
		}
		if err != nil {
			return err
		}
		for _, subDirectory := range subDirectories {
			fmt.Fprintf(out, "%s\\\n", subDirectory)
			if err := Truncate(subDirectory, whenError, out); err != nil {
				return err
			}
		}
	}
	if (attr & windows.FILE_ATTRIBUTE_READONLY) != 0 {
		SetFileAttributes(folder, attr&^windows.FILE_ATTRIBUTE_READONLY)
	}
	if err := windows.Rmdir(folder); err != nil {
		return fmt.Errorf("%s: windows.Rmdir: %w", folder, err)
	}
	return nil
}
