package commands

import "syscall"

import "../interpreter"

func cmd_mkdir(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	for i := 1; i < len(cmd.Args); i++ {
		err := syscall.Mkdir(cmd.Args[i], 0777)
		if err != nil {
			return interpreter.CONTINUE, err
		}
	}
	return interpreter.CONTINUE, nil
}

func cmd_rmdir(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	for i, n := 1, len(cmd.Args); i < n; i++ {
		err := syscall.Rmdir(cmd.Args[i])
		if err != nil {
			return interpreter.CONTINUE, err
		}
	}
	return interpreter.CONTINUE, nil
}
