package commands

import "fmt"
import "os"
import "../interpreter"

func cmd_pwd(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	wd, _ := os.Getwd()
	fmt.Fprintln(cmd.Stdout, wd)
	return interpreter.CONTINUE, nil
}
