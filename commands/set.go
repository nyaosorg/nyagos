package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/nyaosorg/go-readline-ny"

	"github.com/nyaosorg/nyagos/completion"
	"github.com/nyaosorg/nyagos/nodos"
	"github.com/nyaosorg/nyagos/shell"

	"github.com/nyaosorg/nyagos/internal/go-ignorecase-sorted"
)

// ReadStdinAsFile is the flat to read commands from stdin as a file stream
var ReadStdinAsFile = false

type optionT struct {
	V       *bool
	Usage   string
	NoUsage string
}

// BoolOptions are the all global option list.
var BoolOptions = ignoreCaseSorted.MapToDictionary(map[string]*optionT{
	"completion_hidden": {
		V:       &completion.IncludeHidden,
		Usage:   "Include hidden files on completion",
		NoUsage: "Do not include hidden files on completion",
	},
	"completion_slash": {
		V:       &completion.UseSlash,
		Usage:   "use forward slash on completion",
		NoUsage: "Do not use slash on completion",
	},
	"glob": {
		V:       &shell.WildCardExpansionAlways,
		Usage:   "Enable to expand wildcards",
		NoUsage: "Disable to expand wildcards",
	},
	"noclobber": {
		V:       &shell.NoClobber,
		Usage:   "forbide to overwrite files on redirect",
		NoUsage: "Do not forbide to overwrite files no redirect",
	},
	"usesource": {
		V:       &shell.UseSourceRunBatch,
		Usage:   "allow batchfile to change environment variables of nyagos",
		NoUsage: "forbide batchfile to change environment variables of nyagos",
	},
	"tilde_expansion": {
		V:       &shell.TildeExpansion,
		Usage:   "Enable Tilde Expansion",
		NoUsage: "Disable Tilde Expansion",
	},
	"read_stdin_as_file": {
		V:       &ReadStdinAsFile,
		Usage:   "Read commands from stdin as a file stream. Disable to edit line",
		NoUsage: "Read commands from stdin as Windows Console(tty). Enable to edit line",
	},
	"output_surrogate_pair": {
		V:       &readline.SurrogatePairOk,
		Usage:   "Output surrogate pair characters as it is",
		NoUsage: "Output surrogate pair characters like <NNNNN>",
	},
})

func dumpBoolOptions(out io.Writer) {
	max := 0
	for p := BoolOptions.Front(); p != nil; p = p.Next() {
		if L := len(p.Key); L > max {
			max = L
		}
	}
	for p := BoolOptions.Front(); p != nil; p = p.Next() {
		key := p.Key
		val := p.Value
		if *val.V {
			fmt.Fprint(out, "-o ")
		} else {
			fmt.Fprint(out, "+o ")
		}
		fmt.Fprintf(out, "%-*s", max, key)
		if *val.V {
			fmt.Fprintf(out, " (%s)\n", val.Usage)
		} else {
			fmt.Fprintf(out, " (%s)\n", val.NoUsage)
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
				if ptr, ok := BoolOptions.Load(args[0]); ok {
					*ptr.V = true
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
				if ptr, ok := BoolOptions.Load(args[0]); ok {
					*ptr.V = false
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
