package commands

import "fmt"
import "os"
import "os/exec"
import "../interpreter"

func cmd_pwd(cmd *exec.Cmd) (interpreter.NextT, error) {
	wd, _ := os.Getwd()
	fmt.Fprintln(cmd.Stdout, wd)
	return interpreter.CONTINUE, nil
}
