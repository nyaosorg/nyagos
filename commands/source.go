package commands

import (
	"context"
	"github.com/zetamatta/nyagos/shell"
	"io"
)

func cmd_source(ctx context.Context, cmd *shell.Cmd) (int, error) {
	var verbose io.Writer
	args := make([]string, 0, len(cmd.Args))
	debug := false
	for _, arg1 := range cmd.Args[1:] {
		switch arg1 {
		case "-v":
			verbose = cmd.Stderr
		case "-d":
			debug = true
		default:
			args = append(args, arg1)
		}
	}
	if len(cmd.Args) <= 0 {
		return 255, nil
	}

	return shell.Source(args, verbose, debug, cmd.Stdin, cmd.Stdout, cmd.Stderr)
}
