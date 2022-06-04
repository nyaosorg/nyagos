package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/mitchellh/go-ps"
)

func cmdPs(ctx context.Context, cmd Param) (int, error) {
	processes, err := ps.Processes()
	if err != nil {
		return 1, err
	}
	fmt.Fprintf(cmd.Out(), "%6s %6s %s\n", "PID", "PPID", "COMMAND")
	self := os.Getpid()
	for _, p := range processes {
		fmt.Fprintf(cmd.Out(), "%6d %6d %s", p.Pid(), p.PPid(), p.Executable())
		if self == p.Pid() {
			fmt.Fprintln(cmd.Out(), " [self]")
		} else {
			fmt.Fprintln(cmd.Out())
		}
	}
	return 0, nil
}
