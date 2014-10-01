package commands

import "fmt"
import "os/exec"
import "strings"

import "github.com/shiena/ansicolor"

import "../interpreter"

func cmd_echo(cmd *exec.Cmd) (interpreter.NextT, error) {
	fmt.Fprintln(ansicolor.NewAnsiColorWriter(cmd.Stdout), strings.Join(cmd.Args[1:], " "))
	return interpreter.CONTINUE, nil
}
