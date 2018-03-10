package commands

import (
	"context"

	"github.com/zetamatta/nyagos/shell"
)

type paramimp_t struct{ *shell.Cmd }

func Exec(ctx context.Context, cmd *shell.Cmd) (int, bool, error) {
	return exec(ctx, &paramimp_t{cmd})
}
