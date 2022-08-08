package commands

import (
	"context"
	"io"
)

func cmdCls(ctx context.Context, cmd Param) (int, error) {
	io.WriteString(cmd.Term(), "\x1B[1;1H\x1B[2J")
	return 0, nil
}
