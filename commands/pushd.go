package commands

import (
	"errors"
	"fmt"

	"../dos"
	. "../interpreter"
)

var dirstack = make([]string, 0, 20)

func cmd_dirs(cmd *Interpreter) (ErrorLevel, error) {
	wd, _ := dos.Getwd()
	fmt.Fprint(cmd.Stdout, wd)
	for i := len(dirstack) - 1; i >= 0; i-- {
		fmt.Fprint(cmd.Stdout, " ", dirstack[i])
	}
	fmt.Fprintln(cmd.Stdout)
	return CONTINUE, nil
}

func cmd_popd(cmd *Interpreter) (ErrorLevel, error) {
	if len(dirstack) <= 0 {
		return CONTINUE, errors.New("popd: directory stack empty.")
	}
	err := dos.Chdir(dirstack[len(dirstack)-1])
	if err != nil {
		return CONTINUE, err
	}
	dirstack = dirstack[:len(dirstack)-1]
	return cmd_dirs(cmd)
}

func cmd_pushd(cmd *Interpreter) (ErrorLevel, error) {
	wd, err := dos.Getwd()
	if err != nil {
		return CONTINUE, err
	}
	if len(cmd.Args) >= 2 {
		dirstack = append(dirstack, wd)
		err := dos.Chdir(cmd.Args[1])
		if err != nil {
			return CONTINUE, err
		}
	} else {
		if len(dirstack) <= 0 {
			return CONTINUE, errors.New("pushd: directory stack empty.")
		}
		err := dos.Chdir(dirstack[len(dirstack)-1])
		if err != nil {
			return CONTINUE, err
		}
		dirstack[len(dirstack)-1] = wd
	}
	return cmd_dirs(cmd)
}
