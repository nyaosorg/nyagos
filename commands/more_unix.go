// +build !windows

package commands

import (
	"context"
	"errors"
)

func cmdMore(ctx context.Context, cmd Param) (int, error) {
	return 1, errors.New("not support on Linux")
}
