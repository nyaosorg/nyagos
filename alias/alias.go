package alias

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"

	"../interpreter"
)

type Callable interface {
	String() string
	Call(cmd *interpreter.Interpreter) (interpreter.NextT, error)
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

func (this *AliasFunc) Call(cmd *interpreter.Interpreter) (next interpreter.NextT, err error) {
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
		var buffer bytes.Buffer
		buffer.WriteString(this.BaseStr)
		buffer.WriteRune(' ')
		buffer.WriteString(quoteAndJoin(cmd.Args[1:]))
		cmdline = buffer.String()
	}
	it := cmd.Clone()
	it.HookCount = cmd.HookCount + 1
	next, err = it.Interpret(cmdline)
	return
}

var Table = map[string]Callable{}
var paramMatch = regexp.MustCompile("\\$(\\*|[0-9]+)")

func quoteAndJoin(list []string) string {
	var buffer bytes.Buffer
	for _, value := range list {
		if buffer.Len() > 0 {
			buffer.WriteRune(' ')
		}
		buffer.WriteRune('"')
		buffer.WriteString(value)
		buffer.WriteRune('"')
	}
	return buffer.String()
}

var nextHook interpreter.HookT

func hook(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	if cmd.HookCount > 0 {
		return nextHook(cmd)
	}
	callee, ok := Table[strings.ToLower(cmd.Args[0])]
	if !ok {
		return nextHook(cmd)
	}
	if cmd.IsBackGround {
		go func() {
			callee.Call(cmd)
			if cmd.Closer != nil {
				cmd.Closer.Close()
			}
		}()
		return interpreter.CONTINUE, nil
	} else {
		next, err := callee.Call(cmd)
		if next != interpreter.THROUGH && cmd.Closer != nil {
			cmd.Closer.Close()
		}
		return next, err
	}
}

func Init() {
	nextHook = interpreter.SetHook(hook)
}
