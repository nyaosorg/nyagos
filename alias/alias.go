package alias

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/shell"
)

var dbg = false

type callableT interface {
	String() string
	Call(ctx context.Context, cmd *shell.Cmd) (int, error)
}

// Func is the type for string-type alias. It has a Call method
type Func struct {
	BaseStr string
}

// New is the constructor for Func
func New(baseStr string) *Func {
	return &Func{baseStr}
}

// String is the method to support fmt.Stringer
func (f *Func) String() string {
	return f.BaseStr
}

// Call is the method to support callableT and it calls the alias-function.
func (f *Func) Call(ctx context.Context, cmd *shell.Cmd) (next int, err error) {
	isReplaced := false
	if dbg {
		print("Func.Call('", cmd.Arg(0), "')\n")
	}
	cmdline := paramMatch.ReplaceAllStringFunc(f.BaseStr, func(s string) string {
		if s == "$~*" {
			isReplaced = true
			if cmd.Args() != nil && len(cmd.Args()) >= 2 {
				return strings.Join(cmd.Args()[1:], " ")
			}
			return ""
		} else if s == "$*" {
			isReplaced = true
			if cmd.Args() != nil && len(cmd.Args()) >= 2 {
				return strings.Join(cmd.RawArgs()[1:], " ")
			}
			return ""
		} else if len(s) >= 3 && s[0] == '$' && s[1] == '~' && unicode.IsDigit(rune(s[2])) {
			i, err := strconv.ParseInt(s[2:], 10, 0)
			if err == nil {
				isReplaced = true
				if 0 <= i && cmd.Args() != nil && int(i) < len(cmd.Args()) {
					return cmd.Arg(int(i))
				}
				return ""
			}
		}
		i, err := strconv.ParseInt(s[1:], 10, 0)
		if err == nil {
			isReplaced = true
			if 0 <= i && cmd.Args() != nil && int(i) < len(cmd.Args()) {
				return cmd.RawArg(int(i))
			}
			return ""
		}
		return s
	})

	if !isReplaced {
		var buffer strings.Builder
		buffer.WriteString(f.BaseStr)
		for _, s := range cmd.RawArgs()[1:] {
			fmt.Fprintf(&buffer, " %s", s)
		}
		cmdline = buffer.String()
	}
	if dbg {
		print("replaced cmdline=='", cmdline, "'\n")
	}
	next, err = cmd.Interpret(ctx, cmdline)
	return
}

// Table is the ALL ALIAS table !
var Table = map[string]callableT{}
var paramMatch = regexp.MustCompile(`\$(\~)?(\*|[0-9]+)`)

// AllNames returns all-alias names for completion
func AllNames() []completion.Element {
	names := make([]completion.Element, 0, len(Table))
	for name1 := range Table {
		names = append(names, completion.Element1(name1))
	}
	return names
}

var nextHook shell.HookT

type noAliasT string

func hook(ctx context.Context, cmd *shell.Cmd) (int, bool, error) {
	lowerName := strings.ToLower(cmd.Arg(0))
	ctxKey := noAliasT(lowerName)
	if ctx.Value(ctxKey) != nil {
		return nextHook(ctx, cmd)
	}
	callee, ok := Table[lowerName]
	if !ok {
		return nextHook(ctx, cmd)
	}
	next, err := callee.Call(context.WithValue(ctx, ctxKey, true), cmd)
	return next, true, err
}

// Init is the package initializer which inserts hook-function into shell package.
func Init() {
	nextHook = shell.SetHook(hook)
}
