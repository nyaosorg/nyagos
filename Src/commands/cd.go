package commands

import (
	"os"

	"../dos"
	"../interpreter"
)

var prevDir string

func cmd_cd(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	if len(cmd.Args) >= 2 {
		prevDir_, err := os.Getwd()
		if err != nil {
			return interpreter.CONTINUE, err
		}
		if cmd.Args[1] == "-" {
			err = dos.Chdir(prevDir)
		} else {
			err = dos.Chdir(cmd.Args[1])
		}
		prevDir = prevDir_
		return interpreter.CONTINUE, err
	}
	home := dos.GetHome()
	if home != "" {
		prevDir, _ = os.Getwd()
		return interpreter.CONTINUE, dos.Chdir(home)
	}
	return cmd_pwd(cmd)
}
