package commands

import (
	"fmt"
	"os"

	"../interpreter"
)

func cmd_pwd(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	wd, _ := os.Getwd()
	fmt.Fprintln(cmd.Stdout, wd)
	return interpreter.CONTINUE, nil
}
