package alias

import "bytes"
import "os/exec"
import "regexp"
import "strconv"
import "strings"
import "io"

import "../interpreter"

type Callable interface {
	String() string
	Call(cmd *exec.Cmd) (interpreter.NextT, error)
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

var expanded = false

func (this *AliasFunc) Call(cmd *exec.Cmd) (next interpreter.NextT, err error) {
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
	stdio := interpreter.Stdio{
		Stdin:  cmd.Stdin,
		Stdout: cmd.Stdout,
		Stderr: cmd.Stderr,
	}
	expanded = true
	next, err = interpreter.Interpret(cmdline, &stdio)
	expanded = false
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

func hook(cmd *exec.Cmd, IsBackground bool, closer io.Closer) (interpreter.NextT, error) {
	if expanded {
		return nextHook(cmd, IsBackground, closer)
	}
	callee, ok := Table[strings.ToLower(cmd.Args[0])]
	if !ok {
		return nextHook(cmd, IsBackground, closer)
	}
	next, err := callee.Call(cmd)
	if next != interpreter.THROUGH && closer != nil {
		closer.Close()
	}
	return next, err
}

func Init() {
	nextHook = interpreter.SetHook(hook)
}
