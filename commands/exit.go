package commands

import (
	"context"
	"io"
	"os/exec"
)

func cmd_exit(ctx context.Context, cmd *exec.Cmd) (int, error) {
	return 0, io.EOF
}
