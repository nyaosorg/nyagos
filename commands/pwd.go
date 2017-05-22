package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/zetamatta/nyagos/shell"
)

func cmd_pwd(ctx context.Context, cmd *shell.Cmd) (int, error) {
	if len(cmd.Args) >= 2 {
		if i, err := strconv.ParseInt(cmd.Args[1], 10, 0); err == nil && i < 0 {
			i += int64(len(cd_history))
			if i < 0 {
				return NO_HISTORY, fmt.Errorf("pwd %s: too old history", cmd.Args[1])
			}
			fmt.Fprintln(cmd.Stdout, cd_history[i])
			return 0, nil
		}
	} else {
		wd, _ := os.Getwd()
		fmt.Fprintln(cmd.Stdout, wd)
	}
	return 0, nil
}
