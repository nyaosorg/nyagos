package commands

import "os/exec"
import "../interpreter"

func cmd_exit(cmd *exec.Cmd) (interpreter.NextT, error) {
	return interpreter.SHUTDOWN, nil
}
