package commands

import (
	"context"
	"github.com/zetamatta/nyagos/shell"
	"io"
)

func cmdSource(ctx context.Context, cmd Param) (int, error) {
	var verbose io.Writer
	args := make([]string, 0, len(cmd.Args()))
	debug := false
	for _, arg1 := range cmd.Args()[1:] {
		switch arg1 {
		case "-v":
			verbose = cmd.Err()
		case "-d":
			debug = true
		default:
			args = append(args, arg1)
		}
	}
	if len(cmd.Args()) <= 0 {
		return 255, nil
	}

	return shell.Source(args, verbose, debug, cmd.In(), cmd.Out(), cmd.Err())
}
