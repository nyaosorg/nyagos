package commands

import (
	"context"
	"io"
	"regexp"
	"strings"

	"github.com/zetamatta/go-findfile"

	"github.com/nyaosorg/nyagos/completion"
	"github.com/nyaosorg/nyagos/nodos"
	"github.com/nyaosorg/nyagos/shell"
)

// Param is the interface for built-in command
type Param interface {
	Arg(int) string
	Args() []string
	SetArgs(s []string)
	In() io.Reader
	Out() io.Writer
	Err() io.Writer
	Term() io.Writer
	RawArgs() []string
	Spawnlp(context.Context, []string, []string) (int, error)
	Spawnlpe(context.Context, []string, []string, map[string]string) (int, error)
	Loop(context.Context, shell.Stream) (int, error)
	ReadCommand(context.Context) (context.Context, string, error)
	DumpEnv() []string
	Setenv(key, val string)
	GetHistory() shell.History
}

var buildInCommand map[string]func(context.Context, Param) (int, error)
var unscoNamePattern = regexp.MustCompile("^__(.*)__$")
var backslashPattern = regexp.MustCompile(`^\\(\w*)$`)

// Exec is the entry function to call built-in functions from Shell
func Exec(ctx context.Context, cmd Param) (int, bool, error) {
	name := strings.ToLower(cmd.Arg(0))
	if len(name) == 2 && strings.HasSuffix(name, ":") {
		pushCdHistory()
		_, err := nodos.Chdrive(name)
		return 0, true, err
	}
	function, ok := buildInCommand[name]
	if !ok {
		m := unscoNamePattern.FindStringSubmatch(name)
		if m != nil {
			name = m[1]
			function, ok = buildInCommand[name]
			if !ok {
				return 0, false, nil
			}
		} else {
			n := backslashPattern.FindStringSubmatch(name)
			if n == nil {
				return 0, false, nil
			}
			name = n[1]
			function, ok = buildInCommand[name]
			if !ok {
				return 0, false, nil
			}
		}
	}
	cmd.SetArgs(findfile.Globs(cmd.Args()))
	next, err := function(ctx, cmd)
	return next, true, err
}

// AllNames returns all command-names for completion package.
func AllNames(ctx context.Context) ([]completion.Element, error) {
	names := make([]completion.Element, 0, len(buildInCommand))
	for name1 := range buildInCommand {
		names = append(names, completion.Element1(name1))
	}
	return names, nil
}
