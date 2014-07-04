package commands

import "os/exec"

import "../interpreter"

func cmd_rem(cmd *exec.Cmd) (interpreter.NextT, error) {
	return interpreter.CONTINUE, nil
}
