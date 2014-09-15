package commands

import "os/exec"

import "../conio"
import "../interpreter"

func cmd_cls(cmd *exec.Cmd) (interpreter.NextT, error) {
	conio.Cls()
	return interpreter.CONTINUE, nil
}
