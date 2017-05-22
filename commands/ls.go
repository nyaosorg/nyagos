package commands

import (
	"context"
	"io"
	"os"

	"github.com/mattn/go-colorable"

	"github.com/zetamatta/nyagos/commands/ls"
	"github.com/zetamatta/nyagos/shell"
)

func cmd_ls(ctx context.Context, cmd *shell.Cmd) (int, error) {
	var out io.Writer
	if cmd.Stdout == os.Stdout {
		out = colorable.NewColorableStdout()
	} else {
		out = cmd.Stdout
	}
	return 0, ls.Main(ctx, cmd.Args[1:], out, cmd.Stderr)
}
