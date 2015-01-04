package commands

import . "../interpreter"

func cmd_rem(cmd *Interpreter) (NextT, error) {
	return CONTINUE, nil
}
