package commands

import (
	"context"
	"fmt"

	"github.com/mattn/go-colorable"

	"../shell"
)

func cmd_cls(ctx context.Context, cmd *shell.Cmd) (int, error) {
	fmt.Fprint(colorable.NewColorableStdout(), "\x1B[1;1H\x1B[2J")
	return 0, nil
}
