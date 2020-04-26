// +build !windows

package nodos

import "errors"

func shOpenWithDialog(_, _ string) (err error) {
	return errors.New("not supported")
}
