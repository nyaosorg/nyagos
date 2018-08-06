package commands

import (
	"context"
	"fmt"

	"github.com/mitchellh/go-ps"
)

func cmdPs(ctx context.Context, cmd Param) (int, error) {
	processes, err := ps.Processes()
	if err != nil {
		return 1, err
	}
	fmt.Fprintf(cmd.Out(), "%6s %6s %s\n", "PID", "PPID", "COMMAND")
	for _, p := range processes {
		fmt.Fprintf(cmd.Out(), "%6d %6d %s\n", p.Pid(), p.PPid(), p.Executable())
	}
	return 0, nil
}
