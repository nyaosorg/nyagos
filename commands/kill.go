package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

func cmdKill(ctx context.Context, cmd Param) (int, error) {
	args := cmd.Args()
	if len(args) < 2 {
		return 1, fmt.Errorf("Usage: %s PID", args[0])
	}

	pid, err := strconv.Atoi(cmd.Arg(1))
	if err != nil {
		return 1, fmt.Errorf("%s: arguments must be process ID", args[0])
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return 1, err
	}

	err = process.Kill()
	if err != nil {
		return 1, err
	}
	return 0, nil
}
