package commands

import (
	"fmt"
	"strings"

	"github.com/mattn/go-isatty"

	. "../interpreter"
)

func cmd_echo(cmd *Interpreter) (ErrorLevel, error) {
	fmt.Fprint(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	if isatty.IsTerminal(cmd.Stdio[1].Fd()) {
		fmt.Fprint(cmd.Stdout, "\n")
	} else {
		fmt.Fprint(cmd.Stdout, "\r\n")
	}
	return NOERROR, nil
}
