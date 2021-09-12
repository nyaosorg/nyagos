package commands

import (
	"context"
	"fmt"

	"github.com/nyaosorg/go-readline-ny"
)

func cmdBindkey(ctx context.Context, cmd Param) (int, error) {
	if len(cmd.Args()) < 3 {
		fmt.Fprintf(cmd.Err(), "%[1]s: Usage %[1]s KEYNAME FUNCNAME\n",
			cmd.Arg(0))
		return 0, nil
	}
	err := readline.GlobalKeyMap.BindKeySymbol(cmd.Arg(1), cmd.Arg(2))
	if err != nil {
		return 1, err
	}
	return 0, nil
}
