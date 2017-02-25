package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"../cpath"
	"../dos"
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

const (
	CHDIR_FAIL = 1
	NO_HISTORY = 2
)

func cmd_cd_sub(dir string) (int, error) {
	const fileHead = "file:///"

	if strings.HasPrefix(dir, fileHead) {
		dir = dir[len(fileHead):]
	}
	if dir_, err := cpath.CorrectCase(dir); err == nil {
		// println(dir, "->", dir_)
		dir = dir_
	}
	err := dos.Chdir(dir)
	if err == nil {
		return 0, nil
	} else {
		return CHDIR_FAIL, err
	}
}

func cmd_cd(ctx context.Context, cmd *exec.Cmd) (int, error) {
	if len(cmd.Args) >= 2 {
		if cmd.Args[1] == "-" {
			if len(cd_history) < 1 {
				return NO_HISTORY, errors.New("cd - : there is no previous directory")

			}
			directory := cd_history[len(cd_history)-1]
			push_cd_history()
			return cmd_cd_sub(directory)
		} else if cmd.Args[1] == "--history" {
			dir, dir_err := os.Getwd()
			if dir_err == nil {
				fmt.Fprintln(cmd.Stdout, dir)
			} else {
				fmt.Fprintln(cmd.Stderr, dir_err.Error())
			}
			for i := len(cd_history) - 1; i >= 0; i-- {
				fmt.Fprintln(cmd.Stdout, cd_history[i])
			}
			return 0, nil
		} else if cmd.Args[1] == "-h" || cmd.Args[1] == "?" {
			i := len(cd_history) - 10
			if i < 0 {
				i = 0
			}
			for ; i < len(cd_history); i++ {
				fmt.Fprintf(cmd.Stdout, "%d %s\n", i-len(cd_history), cd_history[i])
			}
			return 0, nil
		} else if i, err := strconv.ParseInt(cmd.Args[1], 10, 0); err == nil && i < 0 {
			i += int64(len(cd_history))
			if i < 0 {
				return NO_HISTORY, fmt.Errorf("cd %s: too old history", cmd.Args[1])
			}
			directory := cd_history[i]
			push_cd_history()
			return cmd_cd_sub(directory)
		}
		if strings.EqualFold(cmd.Args[1], "/D") {
			// ignore /D
			cmd.Args = cmd.Args[1:]
		}
		push_cd_history()
		return cmd_cd_sub(strings.Join(cmd.Args[1:], " "))
	}
	home := cpath.GetHome()
	if home != "" {
		push_cd_history()
		return cmd_cd_sub(home)
	}
	return cmd_pwd(ctx, cmd)
}
