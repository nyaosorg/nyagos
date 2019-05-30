package commands

import (
	"context"
	"strings"

	"github.com/zetamatta/nyagos/shell"
)

func cmdExeSc(ctx context.Context, cmd Param) (int, error) {
	return shell.CmdExe{
		Cmdline: strings.Join(cmd.RawArgs()[1:], " "),
		Stdin:   cmd.In(),
		Stdout:  cmd.Out(),
		Stderr:  cmd.Err(),
		Env:     cmd.DumpEnv(),
	}.Run()
}
