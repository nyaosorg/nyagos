// +build !windows

package commands

import (
	"context"
	"errors"
)

func cmdLnk(_ context.Context, cmd1 Param) (int, error) {
	return 1, errors.New("not support on Linux")
}
