package commands

import (
	"context"
	"io"

	"../shell"
)

func cmd_exit(ctx context.Context, cmd *shell.Cmd) (int, error) {
	return 0, io.EOF
}
