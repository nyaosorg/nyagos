package commands

import (
	"fmt"

	"../dos"
	"../interpreter"
)

func cmd_pwd(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	wd, _ := dos.Getwd()
	fmt.Fprintln(cmd.Stdout, wd)
	return interpreter.CONTINUE, nil
}
