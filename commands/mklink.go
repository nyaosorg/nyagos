package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zetamatta/nyagos/nodos"
)

func cmdMklink(_ context.Context, cmd Param) (rc int, err error) {
	args := cmd.Args()
	f := os.Symlink
	label := "Symlink link"
	for len(args) >= 2 && len(args[1]) >= 1 && args[1][0] == '/' {
		switch strings.ToUpper(args[1]) {
		case "/J":
			f = nodos.CreateJunction
			label = "Junction"
		case "/D":
			f = os.Symlink
			label = "Symlink link"
		case "/H":
			f = os.Link
			label = "Hardlink"
		default:
			return 1, fmt.Errorf("Invalid switch - \"%s\".", args[1][1:])
		}
		args = args[1:]
	}
	if len(args) < 3 {
		return 2, fmt.Errorf("The syntax of the command is incorrect")
	}
	err = f(args[2], args[1])
	if err != nil {
		return 3, err
	}
	fmt.Fprintf(cmd.Err(), "%s created for %s <<===>> %s\n", label, args[1], args[2])
	return 0, nil
}
