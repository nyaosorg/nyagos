package commands

import (
	"../conio"
	"os/exec"
)

func cmd_cls(cmd *exec.Cmd) (int, error) {
	conio.Cls()
	return 0, nil
}
