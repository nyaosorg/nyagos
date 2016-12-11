package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func cmd_set(ctx context.Context, cmd *exec.Cmd) (int, error) {
	if len(cmd.Args) <= 1 {
		for _, val := range os.Environ() {
			fmt.Fprintln(cmd.Stdout, val)
		}
		return 0, nil
	}
	for _, arg := range cmd.Args[1:] {
		eqlPos := strings.Index(arg, "=")
		if eqlPos < 0 {
			fmt.Fprintf(cmd.Stdout, "%s=%s\n", arg, os.Getenv(arg))
		} else {
			if eqlPos+1 < len(arg) {
				os.Setenv(arg[:eqlPos], arg[eqlPos+1:])
			} else {
				os.Unsetenv(arg[:eqlPos])
			}
		}
	}
	return 0, nil
}
