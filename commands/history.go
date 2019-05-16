package commands

import (
	"context"
	"github.com/zetamatta/nyagos/history"
)

func cmdHistory(ctx context.Context, args Param) (int, error) {
	return history.CmdHistory(ctx, args, args.GetHistory())
}
