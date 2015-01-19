package commands

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"../dos"
	. "../interpreter"
)

var cd_history = make([]string, 0, 100)
var cd_uniq = map[string]int{}

func push_cd_history() {
	directory, err := os.Getwd()
	if err != nil {
		return
	}
	if i, ok := cd_uniq[directory]; ok {
		for ; i < len(cd_history)-1; i++ {
			cd_history[i] = cd_history[i+1]
			cd_uniq[cd_history[i]] = i
		}
		cd_history[i] = directory
		cd_uniq[directory] = i
	} else {
		cd_uniq[directory] = len(cd_history)
		cd_history = append(cd_history, directory)
	}
}

func cmd_cd(cmd *Interpreter) (NextT, error) {
	var err error
	if len(cmd.Args) >= 2 {
		if cmd.Args[1] == "-" {
			if len(cd_history) < 1 {
				return CONTINUE, errors.New("cd - : there is no previous directory")

			}
			directory := cd_history[len(cd_history)-1]
			push_cd_history()
			err = dos.Chdir(directory)
		} else if cmd.Args[1] == "-h" || cmd.Args[1] == "?" {
			i := len(cd_history) - 10
			if i < 0 {
				i = 0
			}
			for ; i < len(cd_history); i++ {
				fmt.Fprintf(cmd.Stdout, "%d %s\n", i-len(cd_history), cd_history[i])
			}
		} else if i, err := strconv.ParseInt(cmd.Args[1], 10, 0); err == nil && i < 0 {
			i += int64(len(cd_history))
			if i < 0 {
				return CONTINUE, fmt.Errorf("cd %s: too old history", cmd.Args[1])
			}
			directory := cd_history[i]
			push_cd_history()
			err = dos.Chdir(directory)
		} else {
			push_cd_history()
			err = dos.Chdir(cmd.Args[1])
		}
		return CONTINUE, err
	}
	home := dos.GetHome()
	if home != "" {
		push_cd_history()
		return CONTINUE, dos.Chdir(home)
	}
	return cmd_pwd(cmd)
}
