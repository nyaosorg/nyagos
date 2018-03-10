package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/zetamatta/nyagos/dos"
)

func cmdPwd(ctx context.Context, cmd Param) (int, error) {
	physical := true
	if len(cmd.Args()) >= 2 {
		if cmd.Arg(1) == "-P" || cmd.Arg(1) == "-p" {
			physical = true
		} else if cmd.Arg(1) == "-L" || cmd.Arg(1) == "-l" {
			physical = false
		} else if i, err := strconv.ParseInt(cmd.Arg(1), 10, 0); err == nil && i < 0 {
			i += int64(len(cd_history))
			if i < 0 {
				return NO_HISTORY, fmt.Errorf("pwd %s: too old history", cmd.Arg(1))
			}
			fmt.Fprintln(cmd.Out(), cd_history[i])
			return 0, nil
		}
	}
	wd, _ := os.Getwd()
	if physical {
		wd = dos.TruePath(wd)
	}
	fmt.Fprintln(cmd.Out(), wd)
	return 0, nil
}
