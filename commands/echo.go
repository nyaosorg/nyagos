package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"

	. "../interpreter"
)

func cmd_echo(cmd *Interpreter) (ErrorLevel, error) {
	fmt.Fprint(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	if f, f_ok := cmd.Stdout.(*os.File); f_ok && isatty.IsTerminal(f.Fd()) {
		fmt.Fprint(cmd.Stdout, "\n")
	} else {
		fmt.Fprint(cmd.Stdout, "\r\n")
	}
	return NOERROR, nil
}
