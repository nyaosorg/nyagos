package commands

import "../conio"
import "../interpreter"

func cmd_cls(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	conio.Cls()
	return interpreter.CONTINUE, nil
}
