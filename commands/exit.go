package commands

import "../interpreter"

func cmd_exit(cmd *interpreter.Interpreter) (interpreter.ErrorLevel, error) {
	return interpreter.SHUTDOWN, nil
}
