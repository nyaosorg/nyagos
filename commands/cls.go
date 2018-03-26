package commands

import (
	"context"
	"fmt"

	"github.com/mattn/go-colorable"
)

func cmdCls(ctx context.Context, _ Param) (int, error) {
	fmt.Fprint(colorable.NewColorableStdout(), "\x1B[1;1H\x1B[2J")
	return 0, nil
}
