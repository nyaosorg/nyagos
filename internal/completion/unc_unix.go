//go:build !windows
// +build !windows

package completion

import (
	"errors"
)

func uncComplete(str string, force bool) ([]Element, error) {
	return nil, errors.New("not supported")
}
