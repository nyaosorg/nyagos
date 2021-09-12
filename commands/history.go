package commands

import (
	"context"
	"github.com/nyaosorg/nyagos/history"
)

func cmdHistory(ctx context.Context, args Param) (int, error) {
	return history.CmdHistory(ctx, args, args.GetHistory())
}
