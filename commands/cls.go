package commands

import (
	"../conio"
	"../interpreter"
)

func cmd_cls(cmd *interpreter.Interpreter) (interpreter.ErrorLevel, error) {
	conio.Cls()
	return interpreter.NOERROR, nil
}
