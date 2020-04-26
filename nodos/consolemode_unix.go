// +build !windows

package nodos

import (
	"errors"
)

type Handle = uintptr

func changeConsoleMode(console Handle, ops ...ModeOp) (func(), error) {
	return func() {}, errors.New("not supported")
}
