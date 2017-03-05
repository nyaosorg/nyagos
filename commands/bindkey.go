package commands

import (
	"context"
	"fmt"
	"os/exec"

	"../readline"
)

func cmd_bindkey(ctx context.Context, cmd *exec.Cmd) (int, error) {
	if len(cmd.Args) < 3 {
		fmt.Fprintf(cmd.Stderr, "%[1]s: Usage %[1]s KEYNAME FUNCNAME\n",
			cmd.Args[0])
	}
	err := readline.BindKeySymbol(cmd.Args[1], cmd.Args[2])
	if err != nil {
		return 1, err
	}
	return 0, nil
}
