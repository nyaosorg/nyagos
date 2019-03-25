// +build !windows

package nodos

import (
	"os"
)

func CreateJunction(target, mountPt string) error {
	return os.Symlink(target, mountPt)
}
