package option

import "os/exec"
import "strings"

import "../alias"
import . "../alias/table"
import "../commands"
import "../interpreter"

func CommandHooks(cmd *exec.Cmd, IsBackground bool) (interpreter.WhatToDoAfterCmd, error) {
	status, _ := alias.Hook(cmd, IsBackground)
	if status != interpreter.THROUGH {
		return status, nil
	}
	return commands.Exec(cmd, IsBackground)
}

func Parse(getArg func() (string, bool)) {
	for {
		arg, ok := getArg()
		if !ok {
			return
		}
		if arg[0] != '-' {
			continue
		}
		for _, o := range arg[1:] {
			switch o {
			case 'a':
				if equation, equationOk := getArg(); equationOk {
					equationArray := strings.SplitN(equation, "=", 2)
					if len(equationArray) >= 2 {
						Table[strings.ToLower(equationArray[0])] =
							equationArray[1]
					} else {
						delete(Table, strings.ToLower(
							equationArray[0]))
					}
				}
			case 'c', 'k':
				if fname, fnameOk := getArg(); fnameOk {
					interpreter.Interpret(fname, CommandHooks, nil)
				}
				if o == 'c' {
					return
				}
			}
		}
	}
}
