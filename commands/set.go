package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func shrink(values ...string) string {
	hash := make(map[string]bool)

	var buffer bytes.Buffer
	for _, value := range values {
		for _, val1 := range filepath.SplitList(value) {
			val1 = strings.TrimSpace(val1)
			if len(val1) > 0 {
				VAL1 := strings.ToUpper(val1)
				if _, ok := hash[VAL1]; !ok {
					hash[VAL1] = true
					if buffer.Len() > 0 {
						buffer.WriteRune(os.PathListSeparator)
					}
					buffer.WriteString(val1)
				}
			}
		}
	}
	return buffer.String()
}

func cmd_set(ctx context.Context, cmd *exec.Cmd) (int, error) {
	if len(cmd.Args) <= 1 {
		for _, val := range os.Environ() {
			fmt.Fprintln(cmd.Stdout, val)
		}
		return 0, nil
	}
	arg := strings.Join(cmd.Args[1:], " ")
	eqlPos := strings.Index(arg, "=")
	if eqlPos < 0 {
		// set NAME
		fmt.Fprintf(cmd.Stdout, "%s=%s\n", arg, os.Getenv(arg))
	} else if eqlPos >= 3 && arg[eqlPos-1] == '+' {
		// set NAME+=VALUE
		right := arg[eqlPos+1:]
		left := arg[:eqlPos-1]
		os.Setenv(left, shrink(os.Getenv(left), right))
	} else if eqlPos >= 3 && arg[eqlPos-1] == '^' {
		// set NAME^=VALUE
		right := arg[eqlPos+1:]
		left := arg[:eqlPos-1]
		os.Setenv(left, shrink(right, os.Getenv(left)))
	} else if eqlPos+1 < len(arg) {
		// set NAME=VALUE
		os.Setenv(arg[:eqlPos], arg[eqlPos+1:])
	} else {
		// set NAME=
		os.Unsetenv(arg[:eqlPos])
	}
	return 0, nil
}
