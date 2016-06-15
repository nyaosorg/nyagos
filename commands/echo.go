package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-isatty"
)

func cmd_echo(cmd *exec.Cmd) (int, error) {
	fmt.Fprint(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	if f, ok := cmd.Stdout.(*os.File); ok && isatty.IsTerminal(f.Fd()) {
		fmt.Fprint(cmd.Stdout, "\n")
	} else {
		fmt.Fprint(cmd.Stdout, "\r\n")
	}
	return 0, nil
}
