package commands

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/mattn/go-colorable"

	"github.com/zetamatta/nyagos/commands/ls"
)

func cmdLs(ctx context.Context, cmd Param) (int, error) {
	var out io.Writer
	if cmd.Out() == os.Stdout {
		cout := bufio.NewWriter(colorable.NewColorableStdout())
		defer cout.Flush()
		out = cout
	} else {
		out = cmd.Out()
	}
	return 0, ls.Main(ctx, cmd.Args()[1:], out, cmd.Err())
}
