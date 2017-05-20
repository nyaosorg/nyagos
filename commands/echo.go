package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/mattn/go-isatty"

	"../shell"
)

func cmd_echo(ctx context.Context, cmd *shell.Cmd) (int, error) {
	fmt.Fprint(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	if isatty.IsTerminal(cmd.Stdout.Fd()) {
		fmt.Fprint(cmd.Stdout, "\n")
	} else {
		fmt.Fprint(cmd.Stdout, "\r\n")
	}
	return 0, nil
}
