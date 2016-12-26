package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"../interpreter"
)

func cmd_if(ctx context.Context, cmd *exec.Cmd) (int, error) {
	// if "xxx" == "yyy"
	args := cmd.Args
	not := false
	start := 1

	option := map[string]bool{}

	for len(args) >= 2 && strings.HasPrefix(args[1], "/") {
		option[strings.ToLower(args[1])] = true
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
			status = (interpreter.LastErrorLevel <= num)
		}
		start += 2
	}

	if not {
		status = !status
	}
	if status {
		it_ := ctx.Value("interpreter")
		if it, ok := it_.(*interpreter.Interpreter); ok {
			return it.InterpretContext(ctx, strings.Join(it.RawArgs[start:], " "))
		} else {
			fmt.Fprintln(cmd.Stderr, "if: not found sub shell instance")
		}
	}
	return 1, nil
}
