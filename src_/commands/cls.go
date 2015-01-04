package commands

import (
	"../conio"
	"../interpreter"
)

func cmd_cls(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	conio.Cls()
	return interpreter.CONTINUE, nil
}
