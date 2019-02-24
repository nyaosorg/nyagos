// +build !windows

package nodos

import (
	"errors"
	"io"
	"os"
)

func copyFile(src, dst string, isFailIfExists bool) error {
	srcFd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFd.Close()

	if isFailIfExists {
		_, err = os.Stat(dst)
		if err == nil {
			return os.ErrExist
		}
	}
	dstFd, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(dstFd, srcFd)

	if err != nil && err != io.EOF {
		dstFd.Close()
		return err
	}
	if err = dstFd.Close(); err != nil {
		return err
	}
	if fi, err := srcFd.Stat(); err != nil {
		return err
	} else {
		modTime := fi.ModTime()
		if err := os.Chtimes(dst, modTime, modTime); err != nil {
			return err
		}
	}
	return nil
}

func moveFile(src, dst string) error {
	return os.Rename(src, dst)
}

func readShortcut(path string) (string, string, error) {
	return "", "", errors.New("ReadShortcut not support")
}
