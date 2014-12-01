package commands

import (
	"fmt"
	"strings"

	. "../interpreter"
)

func cmd_echo(cmd *Interpreter) (NextT, error) {
	fmt.Fprintln(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	return CONTINUE, nil
}
