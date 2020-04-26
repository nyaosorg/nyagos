// +build !windows

package nodos

import (
	"os"

	"golang.org/x/sys/unix"
)

func getFileAttributes(path string) (uint32, error) {
	var stat unix.Stat_t
	err := unix.Stat(path, &stat)
	if err != nil {
		return 0, err
	}
	return stat.Mode, nil
}

func setFileAttributes(path string, attr uint32) error {
	return unix.Chmod(path, attr)
}

const (
	_REPARSE_POINT = os.ModeSymlink
)
