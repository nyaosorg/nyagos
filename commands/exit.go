package commands

import (
	"os/exec"
)

func cmd_exit(cmd *exec.Cmd) (int, error) {
	return SHUTDOWN, nil
}
