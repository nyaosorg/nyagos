package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/zetamatta/nyagos/dos"
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
	if strings.HasSuffix(strings.ToLower(dir), ".lnk") {
		newdir, _, err := dos.ReadShortcut(dir)
		if err == nil && newdir != "" {
			dir = newdir
		}
	}
	if dir_, err := CorrectCase(dir); err == nil {
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

func cmdCd(ctx context.Context, cmd Param) (int, error) {
	args := cmd.Args()
	if len(args) >= 2 {
		if args[1] == "-" {
			if len(cd_history) < 1 {
				return NO_HISTORY, errors.New("cd - : there is no previous directory")

			}
			directory := cd_history[len(cd_history)-1]
			push_cd_history()
			return cmd_cd_sub(directory)
		} else if args[1] == "--history" {
			dir, dir_err := os.Getwd()
			if dir_err == nil {
				fmt.Fprintln(cmd.Out(), dir)
			} else {
				fmt.Fprintln(cmd.Err(), dir_err.Error())
			}
			for i := len(cd_history) - 1; i >= 0; i-- {
				fmt.Fprintln(cmd.Out(), cd_history[i])
			}
			return 0, nil
		} else if args[1] == "-h" || args[1] == "?" {
			i := len(cd_history) - 10
			if i < 0 {
				i = 0
			}
			for ; i < len(cd_history); i++ {
				fmt.Fprintf(cmd.Out(), "%d %s\n", i-len(cd_history), cd_history[i])
			}
			return 0, nil
		} else if i, err := strconv.ParseInt(args[1], 10, 0); err == nil && i < 0 {
			i += int64(len(cd_history))
			if i < 0 {
				return NO_HISTORY, fmt.Errorf("cd %s: too old history", args[1])
			}
			directory := cd_history[i]
			push_cd_history()
			return cmd_cd_sub(directory)
		}
		if strings.EqualFold(args[1], "/D") {
			// ignore /D
			args = args[1:]
		}
		push_cd_history()
		return cmd_cd_sub(strings.Join(args[1:], " "))
	}
	home := dos.GetHome()
	if home != "" {
		push_cd_history()
		return cmd_cd_sub(home)
	}
	return cmdPwd(ctx, cmd)
}
