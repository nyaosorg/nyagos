package alias

import (
	"regexp"
	"strconv"
	"strings"

	"../conio"
	"../interpreter"
)

type Callable interface {
	String() string
	Call(cmd *interpreter.Interpreter) (interpreter.ErrorLevel, error)
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

func (this *AliasFunc) Call(cmd *interpreter.Interpreter) (next interpreter.ErrorLevel, err error) {
	isReplaced := false
	cmdline := paramMatch.ReplaceAllStringFunc(this.BaseStr, func(s string) string {
		if s == "$*" {
			isReplaced = true
			return quoteAndJoin(cmd.Args[1:])
		}
		i, err := strconv.ParseInt(s[1:], 10, 0)
		if err == nil {
			isReplaced = true
			if 0 <= i && int(i) < len(cmd.Args) {
				return cmd.Args[i]
			}
		}
		return s
	})

	if !isReplaced {
		buffer := make([]byte, 0, 200)
		buffer = append(buffer, this.BaseStr...)
		buffer = append(buffer, ' ')
		buffer = append(buffer, quoteAndJoin(cmd.Args[1:])...)
		cmdline = string(buffer)
	}
	it := cmd.Clone()

	newargs := conio.SplitQ(cmdline)
	if strings.EqualFold(newargs[0], cmd.Args[0]) {
		it.HookCount = 100
	} else {
		it.HookCount = cmd.HookCount + 1
	}
	next, err = it.Interpret(cmdline)
	return
}

var Table = map[string]Callable{}
var paramMatch = regexp.MustCompile("\\$(\\*|[0-9]+)")

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

var nextHook interpreter.HookT

func hook(cmd *interpreter.Interpreter) (interpreter.ErrorLevel, error) {
	if cmd.HookCount > 5 {
		return nextHook(cmd)
	}
	callee, ok := Table[strings.ToLower(cmd.Args[0])]
	if !ok {
		return nextHook(cmd)
	}
	next, err := callee.Call(cmd)
	return next, err
}

func Init() {
	nextHook = interpreter.SetHook(hook)
}
