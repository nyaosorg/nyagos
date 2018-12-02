// +build !windows

package commands

import (
	"errors"
)

func readShortCut(dir string) (string, error) {
	return "", errors.New("not support shortcut")
}
