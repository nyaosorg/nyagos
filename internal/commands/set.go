package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nyaosorg/nyagos/internal/config"
	"github.com/nyaosorg/nyagos/internal/nodos"
)

func cmdSet(ctx context.Context, cmd Param) (int, error) {
	args := cmd.Args()
	if len(args) <= 1 {
		for _, val := range os.Environ() {
			fmt.Fprintln(cmd.Out(), val)
		}
		return 0, nil
	}
	args = args[1:]
	for len(args) > 0 {
		if args[0] == "-o" {
			args = args[1:]
			if len(args) < 1 {
				config.Dump(cmd.Out())
			} else {
				if ptr, ok := config.Bools.Load(args[0]); ok {
					ptr.Set(true)
				} else {
					fmt.Fprintf(cmd.Err(), "-o %s: no such option\n", args[0])
				}
				args = args[1:]
			}
		} else if args[0] == "+o" {
			args = args[1:]
			if len(args) < 1 {
				config.Dump(cmd.Out())
			} else {
				if ptr, ok := config.Bools.Load(args[0]); ok {
					ptr.Set(false)
				} else {
					fmt.Fprintf(cmd.Err(), "+o %s: no such option\n", args[0])
				}
				args = args[1:]
			}
		} else if val := strings.ToLower(args[0]); val == "/a" || val == "-a" {
			value, err := evalEquation(strings.Join(args[1:], " "))
			if err != nil {
				return 1, err
			}
			fmt.Fprintf(cmd.Out(), "%d\n", value)
			return 0, nil
		} else {
			// environment variable operation
			arg := strings.Join(args, " ")
			eqlPos := strings.Index(arg, "=")
			if eqlPos < 0 {
				// set NAME
				fmt.Fprintf(cmd.Out(), "%s=%s\n", arg, os.Getenv(arg))
			} else if eqlPos >= 3 && arg[eqlPos-1] == '+' {
				// set NAME+=VALUE
				right := arg[eqlPos+1:]
				left := arg[:eqlPos-1]
				os.Setenv(left, nodos.JoinList(os.Getenv(left), right))
			} else if eqlPos >= 3 && arg[eqlPos-1] == '^' {
				// set NAME^=VALUE
				right := arg[eqlPos+1:]
				left := arg[:eqlPos-1]
				os.Setenv(left, nodos.JoinList(right, os.Getenv(left)))
			} else if eqlPos+1 < len(arg) {
				// set NAME=VALUE
				os.Setenv(arg[:eqlPos], arg[eqlPos+1:])
			} else {
				// set NAME=
				os.Unsetenv(arg[:eqlPos])
			}
			break
		}
	}
	return 0, nil
}
