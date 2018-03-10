package commands

import (
	"context"
	"fmt"

	"github.com/zetamatta/nyagos/readline"
)

func cmdBindkey(ctx context.Context, cmd Param) (int, error) {
	if len(cmd.Args()) < 3 {
		fmt.Fprintf(cmd.Err(), "%[1]s: Usage %[1]s KEYNAME FUNCNAME\n",
			cmd.Arg(0))
		return 0, nil
	}
	err := readline.BindKeySymbol(cmd.Arg(1), cmd.Arg(2))
	if err != nil {
		return 1, err
	}
	return 0, nil
}
