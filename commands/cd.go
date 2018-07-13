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

var cdHistory = make([]string, 0, 100)
var cdUniq = map[string]int{}

func pushCdHistory() {
	directory, err := os.Getwd()
	if err != nil {
		return
	}
	if i, ok := cdUniq[directory]; ok {
		for ; i < len(cdHistory)-1; i++ {
			cdHistory[i] = cdHistory[i+1]
			cdUniq[cdHistory[i]] = i
		}
		cdHistory[i] = directory
		cdUniq[directory] = i
	} else {
		cdUniq[directory] = len(cdHistory)
		cdHistory = append(cdHistory, directory)
	}
}

const (
	errnoChdirFail = 1
	errnoNoHistory = 2
)

func cmdCdSub(dir string) (int, error) {
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
	if dirTmp, err := CorrectCase(dir); err == nil {
		// println(dir, "->", dirTmp)
		dir = dirTmp
	}
	err := dos.Chdir(dir)
	if err == nil {
		return 0, nil
	}
	return errnoChdirFail, err
}

func cmdCd(ctx context.Context, cmd Param) (int, error) {
	args := cmd.Args()
	if len(args) >= 2 {
		if args[1] == "-" {
			if len(cdHistory) < 1 {
				return errnoNoHistory, errors.New("cd - : there is no previous directory")

			}
			directory := cdHistory[len(cdHistory)-1]
			pushCdHistory()
			return cmdCdSub(directory)
		} else if args[1] == "--history" {
			dir, err := os.Getwd()
			if err == nil {
				fmt.Fprintln(cmd.Out(), dir)
			} else {
				fmt.Fprintln(cmd.Err(), err.Error())
			}
			for i := len(cdHistory) - 1; i >= 0; i-- {
				fmt.Fprintln(cmd.Out(), cdHistory[i])
			}
			return 0, nil
		} else if args[1] == "-h" || args[1] == "?" {
			i := len(cdHistory) - 10
			if i < 0 {
				i = 0
			}
			for ; i < len(cdHistory); i++ {
				fmt.Fprintf(cmd.Out(), "cd %d => cd \"%s\"\n", i-len(cdHistory), cdHistory[i])
			}
			return 0, nil
		} else if i, err := strconv.ParseInt(args[1], 10, 0); err == nil && i < 0 {
			i += int64(len(cdHistory))
			if i < 0 {
				return errnoNoHistory, fmt.Errorf("cd %s: too old history", args[1])
			}
			directory := cdHistory[i]
			pushCdHistory()
			return cmdCdSub(directory)
		}
		if strings.EqualFold(args[1], "/D") {
			// ignore /D
			args = args[1:]
		}
		pushCdHistory()
		return cmdCdSub(strings.Join(args[1:], " "))
	}
	home := dos.GetHome()
	if home != "" {
		pushCdHistory()
		return cmdCdSub(home)
	}
	return cmdPwd(ctx, cmd)
}
