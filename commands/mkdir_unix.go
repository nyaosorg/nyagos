// +build !windows

package commands

import (
	"context"
	"errors"
)

func cmdMkdir(ctx context.Context, cmd Param) (int, error) {
	return 1, errors.New("not support on Linux")
}

func cmdRmdir(ctx context.Context, cmd Param) (int, error) {
	return 1, errors.New("not support on Linux")
}
