package commands

import (
	"fmt"
	"strconv"

	"../dos"
	. "../interpreter"
)

func cmd_pwd(cmd *Interpreter) (NextT, error) {
	if len(cmd.Args) >= 2 {
		if i, err := strconv.ParseInt(cmd.Args[1], 10, 0); err == nil && i < 0 {
			i += int64(len(cd_history))
			if i < 0 {
				return CONTINUE, fmt.Errorf("pwd %s: too old history", cmd.Args[1])
			}
			fmt.Fprintln(cmd.Stdout, cd_history[i])
			return CONTINUE, nil
		}
	} else {
		wd, _ := dos.Getwd()
		fmt.Fprintln(cmd.Stdout, wd)
	}
	return CONTINUE, nil
}
