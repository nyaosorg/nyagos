package commands

import "fmt"
import "os/exec"
import "strings"

import "../interpreter"

func cmd_echo(cmd *exec.Cmd) (interpreter.NextT, error) {
	fmt.Fprintln(cmd.Stdout, strings.Join(cmd.Args[1:], " "))
	return interpreter.CONTINUE, nil
}
