package option

import "io"
import "os/exec"

import "../alias"
import "../commands"
import "../interpreter"

func CommandHooks(cmd *exec.Cmd, IsBackground bool, closer io.Closer) (interpreter.NextT, error) {
	status, _ := alias.Hook(cmd, IsBackground, closer)
	if status != interpreter.THROUGH {
		return status, nil
	}
	return commands.Exec(cmd, IsBackground, closer)
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
