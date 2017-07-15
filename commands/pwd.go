package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/shell"
)

func cmd_pwd(ctx context.Context, cmd *shell.Cmd) (int, error) {
	physical := true
	if len(cmd.Args) >= 2 {
		if cmd.Args[1] == "-P" || cmd.Args[1] == "-p" {
			physical = true
		} else if cmd.Args[1] == "-L" || cmd.Args[1] == "-l" {
			physical = false
		} else if i, err := strconv.ParseInt(cmd.Args[1], 10, 0); err == nil && i < 0 {
			i += int64(len(cd_history))
			if i < 0 {
				return NO_HISTORY, fmt.Errorf("pwd %s: too old history", cmd.Args[1])
			}
			fmt.Fprintln(cmd.Stdout, cd_history[i])
			return 0, nil
		}
	}
	wd, _ := os.Getwd()
	if physical {
		wd = dos.TruePath(wd)
	}
	fmt.Fprintln(cmd.Stdout, wd)
	return 0, nil
}
