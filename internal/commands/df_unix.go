//go:build !windows
// +build !windows

package commands

import (
	"context"
	"errors"
)

func cmdDiskFree(_ context.Context, cmd Param) (int, error) {
	return 1, errors.New("not supported")
}
