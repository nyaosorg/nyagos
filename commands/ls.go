package commands

import "os/exec"

import "github.com/shiena/ansicolor"

import "../interpreter"
import "./ls"

func cmd_ls(cmd *exec.Cmd) (interpreter.NextT, error) {
	return interpreter.CONTINUE,
		ls.Main(cmd.Args[1:], ansicolor.NewAnsiColorWriter(cmd.Stdout))
}
