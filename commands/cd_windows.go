// +build windows

package commands

import (
	"github.com/zetamatta/nyagos/dos"
)

func readShortCut(dir string) (string, error) {
	newdir, _, err := dos.ReadShortcut(dir)
	return newdir, err
}
