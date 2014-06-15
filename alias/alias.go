package alias

import "bytes"
import "os/exec"
import "regexp"
import "strconv"
import "strings"

import "../commands"
import "../interpreter"

var Table = map[string]string{
	"assoc":  "%COMSPEC% /c assoc",
	"attrib": "%COMSPEC% /c attrib",
	"copy":   "%COMSPEC% /c copy",
	"del":    "%COMSPEC% /c del",
	"dir":    "%COMSPEC% /c dir",
	"for":    "%COMSPEC% /c for",
	"md":     "%COMSPEC% /c md",
	"mkdir":  "%COMSPEC% /c mkdir",
	"mklink": "%COMSPEC% /c mklink",
	"move":   "%COMSPEC% /c move",
	"open":   "%COMSPEC% /c start",
	"rd":     "%COMSPEC% /c rd",
	"ren":    "%COMSPEC% /c ren",
	"rename": "%COMSPEC% /c rename",
	"rmdir":  "%COMSPEC% /c rmdir",
	"start":  "%COMSPEC% /c start",
	"type":   "%COMSPEC% /c type",
}
var paramMatch = regexp.MustCompile("\\$(\\*|[0-9]+)")

func Hook(cmd *exec.Cmd, IsBackground bool) (interpreter.WhatToDoAfterCmd, error) {
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
	return interpreter.Interpret(
		cmdline,
		commands.Exec,
		&stdio)
}
