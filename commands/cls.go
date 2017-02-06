package commands

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/mattn/go-colorable"
)

func cmd_cls(ctx context.Context, cmd *exec.Cmd) (int, error) {
	fmt.Fprint(colorable.NewColorableStdout(), "\x1B[1;1H\x1B[2J")
	return 0, nil
}
