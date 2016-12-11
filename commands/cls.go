package commands

import (
	"context"
	"os/exec"

	"../conio"
)

func cmd_cls(ctx context.Context, cmd *exec.Cmd) (int, error) {
	conio.Cls()
	return 0, nil
}
