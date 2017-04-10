package commands

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	"../shell"
)

func cmd_if(ctx context.Context, cmd *shell.Cmd) (int, error) {
	// if "xxx" == "yyy"
	args := cmd.Args
	not := false
	start := 1

	option := map[string]struct{}{}

	for len(args) >= 2 && strings.HasPrefix(args[1], "/") {
		option[strings.ToLower(args[1])] = struct{}{}
		args = args[1:]
		start++
	}

	if len(args) >= 2 && strings.EqualFold(args[1], "not") {
		not = true
		args = args[1:]
		start++
	}
	status := false
	if len(args) >= 4 && args[2] == "==" {
		if _, ok := option["/i"]; ok {
			status = strings.EqualFold(args[1], args[3])
		} else {
			status = (args[1] == args[3])
		}
		args = args[4:]
		start += 3
	} else if len(args) >= 3 && strings.EqualFold(args[1], "exist") {
		_, err := os.Stat(args[2])
		status = (err == nil)
		args = args[3:]
		start += 2
	} else if len(args) >= 3 && strings.EqualFold(args[1], "errorlevel") {
		num, num_err := strconv.Atoi(args[2])
		if num_err == nil {
			lastErrorLevel, ok := ctx.Value("errorlevel").(int)
			if !ok {
				return -1, errors.New("if: could not get context.Value(\"errorlevel\")")
			}
			status = (lastErrorLevel <= num)
		}
		start += 2
	}

	if not {
		status = !status
	}
	if status {
		exec, exec_ok := ctx.Value("exec").(func(string) (int, error))
		if !exec_ok {
			return -1, errors.New("if: could not get context.Value(\"exec\")")
		}
		rawargs, rawargs_ok := ctx.Value("rawargs").([]string)
		if !rawargs_ok {
			return -1, errors.New("if: could not get context.Value(\"rawargs\")")
		}
		cmdline := strings.Join(rawargs[start:], " ")
		// println(cmdline)
		return exec(cmdline)
	} else {
		gotoeol, ok := ctx.Value("gotoeol").(func())
		if !ok {
			return -1, errors.New("if: could not get context.Value(\"gotoeol\")")
		}
		gotoeol()
		return 0, nil
	}
}
