package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

func shrink(values ...string) string {
	hash := make(map[string]struct{})

	var buffer strings.Builder
	for _, value := range values {
		for _, val1 := range filepath.SplitList(value) {
			val1 = strings.TrimSpace(val1)
			if len(val1) > 0 {
				VAL1 := strings.ToUpper(val1)
				if _, ok := hash[VAL1]; !ok {
					hash[VAL1] = struct{}{}
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

var BoolOptions = map[string]*bool{
	"glob":           &shell.WildCardExpansionAlways,
	"noclobber":      &shell.NoClobber,
	"usesource":      &shell.UseSourceRunBatch,
	"cleanup_buffer": &readline.FlushBeforeReadline,
}

func dumpBoolOptions(out io.Writer) {
	for key, val := range BoolOptions {
		fmt.Fprintf(out, "%-16s", key)
		if *val {
			fmt.Fprintln(out, "on")
		} else {
			fmt.Fprintln(out, "off")
		}
	}
}

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
				dumpBoolOptions(cmd.Out())
			} else {
				if ptr, ok := BoolOptions[args[0]]; ok {
					*ptr = true
				} else {
					fmt.Fprintf(cmd.Err(), "-o %s: no such option\n", args[0])
				}
				args = args[1:]
			}
		} else if args[0] == "+o" {
			args = args[1:]
			if len(args) < 1 {
				dumpBoolOptions(cmd.Out())
			} else {
				if ptr, ok := BoolOptions[args[0]]; ok {
					*ptr = false
				} else {
					fmt.Fprintf(cmd.Err(), "+o %s: no such option\n", args[0])
				}
				args = args[1:]
			}
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
			break
		}
	}
	return 0, nil
}
