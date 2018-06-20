package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
	"github.com/zetamatta/nyagos/texts"
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

type optionT struct {
	V       *bool
	Usage   string
	NoUsage string
}

// BoolOptions are the all global option list.
var BoolOptions = map[string]*optionT{
	"cleanup_buffer": {
		V:       &readline.FlushBeforeReadline,
		Usage:   "Clean up key buffer at prompt",
		NoUsage: "Do not clean up key buffer at prompt",
	},
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
}

func dumpBoolOptions(out io.Writer) {
	max := 0
	for key := range BoolOptions {
		if L := len(key); L > max {
			max = L
		}
	}
	for _, key := range texts.SortedKeys(BoolOptions) {
		val := BoolOptions[key]
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
				if ptr, ok := BoolOptions[args[0]]; ok {
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
				if ptr, ok := BoolOptions[args[0]]; ok {
					*ptr.V = false
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
