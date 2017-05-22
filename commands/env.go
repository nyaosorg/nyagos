package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/zetamatta/nyagos/shell"
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

func cmd_env(ctx context.Context, cmd *shell.Cmd) (int, error) {
	args, hash := array2hash(cmd.Args[1:])
	if len(args) <= 0 {
		for _, val := range os.Environ() {
			fmt.Fprintln(cmd.Stdout, val)
		}
		return 0, nil
	}
	backup := map[string]string{}
	for key, val := range hash {
		backup[key] = os.Getenv(key)
		os.Setenv(key, val)
	}

	var rc int
	subCmd, err := cmd.Clone()
	if err == nil {
		subCmd.Args = args
		rc, err = subCmd.SpawnvpContext(ctx)
	} else {
		rc = -1
	}

	for key, val := range backup {
		os.Setenv(key, val)
	}
	return rc, err
}
