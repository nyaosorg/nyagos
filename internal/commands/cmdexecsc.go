package commands

import (
	"context"
	"strings"

	"github.com/nyaosorg/nyagos/internal/source"
)

func cmdExeSc(ctx context.Context, cmd Param) (int, error) {
	return source.System{
		Cmdline: strings.Join(cmd.RawArgs()[1:], " "),
		Stdin:   cmd.In(),
		Stdout:  cmd.Out(),
		Stderr:  cmd.Err(),
		Env:     cmd.DumpEnv(),
	}.Run()
}
