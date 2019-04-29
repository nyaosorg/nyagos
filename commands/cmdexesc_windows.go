package commands

import (
	"context"
	"strings"

	"github.com/zetamatta/nyagos/shell"
)

func cmdExeSc(ctx context.Context, cmd Param) (int, error) {
	return shell.CmdExe(
		strings.Join(cmd.RawArgs()[1:], " "),
		cmd.In(),
		cmd.Out(),
		cmd.Err(),
		cmd.DumpEnv())
}
