package commands

import "fmt"
import "os"
import "os/exec"
import "strings"

import "../interpreter"

func cmd_set(cmd *exec.Cmd) (interpreter.NextT, error) {
	if len(cmd.Args) <= 1 {
		for _, val := range os.Environ() {
			fmt.Fprintln(cmd.Stdout, val)
		}
		return interpreter.CONTINUE, nil
	}
	for _, arg := range cmd.Args[1:] {
		eqlPos := strings.Index(arg, "=")
		if eqlPos < 0 {
			fmt.Fprintf(cmd.Stdout, "%s=%s\n", arg, os.Getenv(arg))
		} else {
			os.Setenv(arg[:eqlPos], arg[eqlPos+1:])
		}
	}
	return interpreter.CONTINUE, nil
}
