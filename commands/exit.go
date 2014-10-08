package commands

import "../interpreter"

func cmd_exit(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	return interpreter.SHUTDOWN, nil
}
