package functions

import (
	"fmt"
	"os/exec"

	"github.com/mattn/msgbox"
	"github.com/nyaosorg/go-windows-mbcs"
	"github.com/nyaosorg/go-windows-su"
)

func CmdMsgBox(args []anyT) []anyT {
	var message string
	title := "nyagos"
	if len(args) >= 1 {
		message = fmt.Sprint(args[0])
	}
	if len(args) >= 2 {
		title = fmt.Sprint(args[1])
	}
	msgbox.Show(0, message, title, msgbox.OK)
	return []anyT{}
}

func CmdElevated([]anyT) []anyT {
	flag, _ := su.IsElevated()
	return []anyT{flag}
}

func CmdShellExecute(args []anyT) []anyT {
	pid, err := su.ShellExecute(
		toStr(args, 0),
		toStr(args, 1),
		toStr(args, 2),
		toStr(args, 3))
	if err != nil {
		return []anyT{nil, err}
	}
	return []anyT{pid}
}

func CmdRawExec(this *Param) []anyT {
	argv := stackToSlice(this)
	xcmd := exec.Command(argv[0], argv[1:]...)
	xcmd.Stdin = this.In
	xcmd.Stdout = this.Out
	xcmd.Stderr = this.Err
	err := xcmd.Run()
	errorlevel := xcmd.ProcessState.ExitCode()
	if err != nil {
		fmt.Fprintln(xcmd.Stderr, err.Error())
		return []anyT{errorlevel, err.Error()}
	}
	return []anyT{errorlevel}
}

func CmdAtoU(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{nil, TooFewArguments}
	}
	if s, ok := args[0].(string); ok {
		val, err := mbcs.AtoU([]byte(s), mbcs.ConsoleCP())
		if err != nil {
			return []anyT{nil, err}
		}
		return []anyT{val}
	}
	return []anyT{fmt.Sprint(args[0])}
}

func CmdUtoA(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{nil, TooFewArguments}
	}
	utf8 := fmt.Sprint(args[0])
	bin, err := mbcs.UtoA(utf8, mbcs.ConsoleCP())
	if err != nil {
		return []anyT{nil, err}
	}
	return []anyT{bin, nil}
}
