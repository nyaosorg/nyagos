// +build !windows

package commands

import (
	"context"
	"errors"
)

func cmdMove(ctx context.Context, cmd Param) (int, error) {
	return 1, errors.New("not support on Linux")
}
func cmdCopy(ctx context.Context, cmd Param) (int, error) {
	return 1, errors.New("not support on Linux")
}
func cmdLn(ctx context.Context, cmd Param) (int, error) {
	return 1, errors.New("not support on Linux")
}
