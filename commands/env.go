package commands

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
)

func array2hash(args []string) ([]string, map[string]string) {
	hash := map[string]string{}
	for i, arg1 := range args {
		equalPos := strings.IndexRune(arg1, '=')
		if equalPos < 0 {
			return args[i:], hash
		}
		key := arg1[:equalPos]
		val := arg1[equalPos+1:]
		hash[key] = val
	}
	return []string{}, hash
}

func cmd_env(ctx context.Context, cmd *exec.Cmd) (int, error) {
	args, hash := array2hash(cmd.Args[1:])
	if len(args) <= 0 {
		return 0, nil
	}
	backup := map[string]string{}
	for key, val := range hash {
		backup[key] = os.Getenv(key)
		os.Setenv(key, val)
	}
	rawargs, ok := ctx.Value("rawargs").([]string)
	if !ok {
		return 0, errors.New("can not get rawargs")
	}
	cmdline := strings.Join(rawargs[(len(rawargs)-len(args)):], " ")
	shell, ok := ctx.Value("exec").(func(string) (int, error))
	if !ok {
		return 0, errors.New("can not shell")
	}
	println("cmdline=" + cmdline)
	rc, err := shell(cmdline)
	for key, val := range backup {
		os.Setenv(key, val)
	}
	return rc, err
}
