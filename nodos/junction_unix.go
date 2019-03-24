// +build !windows

package nodos

import (
	"os"
)

func CreateJunction(mountPt, target string) error {
	return os.Symlink(target, mountPt)
}
