// +build !windows

package commands

import (
	"context"
	"errors"
)

func cmdAttrib(ctx context.Context, cmd Param) (int, error) {
	return 1, errors.New("not supported")
}
