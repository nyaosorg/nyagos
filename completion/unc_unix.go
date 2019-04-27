// +build !windows

package completion

import (
	"errors"
)

func uncComplete(str string) ([]Element, error) {
	return nil, errors.New("not supported")
}
