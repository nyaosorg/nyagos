package commands

import (
	"fmt"
	"strings"

	. "../interpreter"
)

func cmd_echo(cmd *Interpreter) (ErrorLevel, error) {
	fmt.Fprintln(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	return NOERROR, nil
}
