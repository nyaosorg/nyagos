package alias

import "bytes"
import "os/exec"
import "regexp"
import "strconv"
import "strings"
import "io"

import . "./table"
import "../commands"
import "../interpreter"

var paramMatch = regexp.MustCompile("\\$(\\*|[0-9]+)")

func Hook(cmd *exec.Cmd, IsBackground bool, closer io.Closer) (interpreter.NextT, error) {
	baseStr, ok := Table[strings.ToLower(cmd.Args[0])]
	if !ok {
		return interpreter.THROUGH, nil
	}
	isReplaced := false
	cmdline := paramMatch.ReplaceAllStringFunc(baseStr, func(s string) string {
		if s == "$*" {
			isReplaced = true
			return strings.Join(cmd.Args[1:], " ")
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
		buffer.WriteString(baseStr)

		for _, arg := range cmd.Args[1:] {
			buffer.WriteRune(' ')
			buffer.WriteString(arg)
		}

		cmdline = buffer.String()
	}
	var stdio interpreter.Stdio
	stdio.Stdin = cmd.Stdin
	stdio.Stdout = cmd.Stdout
	stdio.Stderr = cmd.Stderr
	nextT, err := interpreter.Interpret(
		cmdline,
		commands.Exec,
		&stdio)
	if nextT != interpreter.THROUGH && closer != nil {
		closer.Close()
	}
	return nextT, err
}
