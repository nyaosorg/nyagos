package commands

import (
	"io"
	"os/exec"
)

func cmd_exit(cmd *exec.Cmd) (int, error) {
	return 0, io.EOF
}
