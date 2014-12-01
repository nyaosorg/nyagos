package commands

import (
	"fmt"
	"os/exec"

	. "../interpreter"
)

func cmd_which(cmd *Interpreter) (NextT, error) {
	for _, name := range cmd.Args[1:] {
		path, err := exec.LookPath(name)
		if err != nil {
			return CONTINUE, err
		}
		fmt.Fprintln(cmd.Stdout, path)
	}
	return CONTINUE, nil
}
