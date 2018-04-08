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

type AliasFunc struct {
	BaseStr string
}

func New(baseStr string) *AliasFunc {
	return &AliasFunc{baseStr}
}

func (this *AliasFunc) String() string {
	return this.BaseStr
}

func (this *AliasFunc) Call(ctx context.Context, cmd *shell.Cmd) (next int, err error) {
	isReplaced := false
	if dbg {
		print("AliasFunc.Call('", cmd.Arg(0), "')\n")
	}
	cmdline := paramMatch.ReplaceAllStringFunc(this.BaseStr, func(s string) string {
		if s == "$~*" {
			isReplaced = true
			if cmd.Args() != nil && len(cmd.Args()) >= 2 {
				return strings.Join(cmd.Args()[1:], " ")
			} else {
				return ""
			}
		} else if s == "$*" {
			isReplaced = true
			if cmd.Args() != nil && len(cmd.Args()) >= 2 {
				return strings.Join(cmd.RawArgs()[1:], " ")
			} else {
				return ""
			}
		} else if len(s) >= 3 && s[0] == '$' && s[1] == '~' && unicode.IsDigit(rune(s[2])) {
			i, err := strconv.ParseInt(s[2:], 10, 0)
			if err == nil {
				isReplaced = true
				if 0 <= i && cmd.Args() != nil && int(i) < len(cmd.Args()) {
					return cmd.Arg(int(i))
				} else {
					return ""
				}
			}
		}
		i, err := strconv.ParseInt(s[1:], 10, 0)
		if err == nil {
			isReplaced = true
			if 0 <= i && cmd.Args() != nil && int(i) < len(cmd.Args()) {
				return cmd.RawArg(int(i))
			} else {
				return ""
			}
		}
		return s
	})

	if !isReplaced {
		var buffer strings.Builder
		buffer.WriteString(this.BaseStr)
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

var Table = map[string]callableT{}
var paramMatch = regexp.MustCompile(`\$(\~)?(\*|[0-9]+)`)

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

func Init() {
	nextHook = shell.SetHook(hook)
}
