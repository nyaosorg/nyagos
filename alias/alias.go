package alias

import "bytes"
import "os/exec"
import "strings"

import "../interpreter"
import "../builtincmd"

var table = map[string]string{
	"l": "ls -oF",
}

func Hook(cmd *exec.Cmd, IsBackground bool) (interpreter.WhatToDoAfterCmd, error) {
	baseStr, ok := table[strings.ToLower(cmd.Args[0])]
	if !ok {
		return interpreter.THROUGH, nil
	}
	var buffer bytes.Buffer
	buffer.WriteString(baseStr)
	for _, arg := range cmd.Args[1:] {
		buffer.WriteRune(' ')
		buffer.WriteString(arg)
	}
	var stdio interpreter.Stdio
	stdio.Stdin = cmd.Stdin
	stdio.Stdout = cmd.Stdout
	stdio.Stderr = cmd.Stderr
	return interpreter.Interpret(
		buffer.String(),
		builtincmd.Exec,
		&stdio)
}
