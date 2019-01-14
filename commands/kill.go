package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
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

func cmdKillAll(ctx context.Context, cmd Param) (int, error) {
	args := cmd.Args()
	if len(args) < 2 {
		return 1, fmt.Errorf("Usage: %s {ExecutableName...}", args[0])
	}
	processes, err := ps.Processes()
	if err != nil {
		return 1, err
	}
	keywords := make([]string, 0, len(args)-1)
	for _, w := range args[1:] {
		if len(w) <= 1 {
			fmt.Fprintf(cmd.Err(), "%s: ExecutableName must be more than 1-character\n", w)
			continue
		}
		keywords = append(keywords, strings.ToUpper(w))
	}
	errorlevel := 1
	for _, p := range processes {
		name := strings.ToUpper(p.Executable())
		for _, w := range keywords {
			if strings.Contains(name, w) {
				process1, err := os.FindProcess(p.Pid())
				if err == nil {
					err = process1.Kill()
				}
				if err == nil {
					fmt.Fprintf(cmd.Out(), "Killed [%d] %s.\n", p.Pid(), p.Executable())
					if errorlevel == 1 {
						errorlevel = 0
					}
				} else {
					fmt.Fprintf(cmd.Err(), "Fail to kill [%d] %s because %s.\n", p.Pid(), p.Executable(), err.Error())
					errorlevel = 2
				}
				break
			}
		}
	}
	return errorlevel, nil
}
