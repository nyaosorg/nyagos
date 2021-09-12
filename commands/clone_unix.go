//go:build !windows
// +build !windows

package commands

import (
	"context"
	"errors"
)

func cmdSu(ctx context.Context, cmd Param) (int, error) {
	return 1, errors.New("not supported")
}
func cmdClone(ctx context.Context, cmd Param) (int, error) {
	return 1, errors.New("not supported")
}
