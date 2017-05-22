package alias

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/shell"
)

var dbg = false

type Callable interface {
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
		print("AliasFunc.Call('", cmd.Args[0], "')\n")
	}
	cmdline := paramMatch.ReplaceAllStringFunc(this.BaseStr, func(s string) string {
		if s == "$~*" {
			isReplaced = true
			if cmd.Args != nil && len(cmd.Args) >= 2 {
				return strings.Join(cmd.Args[1:], " ")
			} else {
				return ""
			}
		} else if s == "$*" {
			isReplaced = true
			if cmd.Args != nil && len(cmd.Args) >= 2 {
				return strings.Join(cmd.RawArgs[1:], " ")
			} else {
				return ""
			}
		} else if len(s) >= 3 && s[0] == '$' && s[1] == '~' && strings.IndexByte("0123456789", s[2]) >= 0 {
			i, err := strconv.ParseInt(s[2:], 10, 0)
			if err == nil {
				isReplaced = true
				if 0 <= i && cmd.Args != nil && int(i) < len(cmd.Args) {
					return cmd.Args[i]
				} else {
					return ""
				}
			}
		}
		i, err := strconv.ParseInt(s[1:], 10, 0)
		if err == nil {
			isReplaced = true
			if 0 <= i && cmd.Args != nil && int(i) < len(cmd.Args) {
				return cmd.RawArgs[i]
			} else {
				return ""
			}
		}
		return s
	})

	if !isReplaced {
		buffer := make([]byte, 0, 200)
		buffer = append(buffer, this.BaseStr...)
		for _, s := range cmd.RawArgs[1:] {
			buffer = append(buffer, ' ')
			buffer = append(buffer, s...)
		}
		cmdline = string(buffer)
	}
	if dbg {
		print("replaced cmdline=='", cmdline, "'\n")
		print("cmd.Clone\n")
	}
	it, err := cmd.Clone()
	if err != nil {
		return 255, err
	}
	if dbg {
		print("done cmd.Clone\n")
	}

	arg1 := shell.QuotedFirstWord(cmdline)
	if strings.EqualFold(arg1, cmd.Args[0]) {
		it.HookCount = 100
	} else {
		it.HookCount = cmd.HookCount + 1
	}
	if dbg {
		print("it.Interpret\n")
	}
	next, err = it.InterpretContext(ctx, cmdline)
	if dbg {
		print("done it.Interpret\n")
	}
	return
}

var Table = map[string]Callable{}
var paramMatch = regexp.MustCompile(`\$(\~)?(\*|[0-9]+)`)

func AllNames() []completion.Element {
	names := make([]completion.Element, 0, len(Table))
	for name1, _ := range Table {
		names = append(names, completion.Element{name1, name1})
	}
	return names
}

func quoteAndJoin(list []string) string {
	buffer := make([]byte, 0, 100)
	for i, value := range list {
		if i > 0 {
			buffer = append(buffer, ' ')
		}
		if strings.IndexByte(value, ' ') >= 0 {
			buffer = append(buffer, '"')
			buffer = append(buffer, value...)
			buffer = append(buffer, '"')
		} else {
			buffer = append(buffer, value...)
		}
	}
	return string(buffer)
}

var nextHook shell.HookT

func hook(ctx context.Context, cmd *shell.Cmd) (int, bool, error) {
	if cmd.HookCount > 5 {
		return nextHook(ctx, cmd)
	}
	callee, ok := Table[strings.ToLower(cmd.Args[0])]
	if !ok {
		return nextHook(ctx, cmd)
	}
	next, err := callee.Call(ctx, cmd)
	return next, true, err
}

func Init() {
	nextHook = shell.SetHook(hook)
}
