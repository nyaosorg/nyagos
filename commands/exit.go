package commands

import (
	"context"
	"io"
)

func cmdExit(_ context.Context, _ Param) (int, error) {
	return 0, io.EOF
}
