package functions

import (
	"fmt"
	"os/exec"

	"github.com/mattn/msgbox"
	"github.com/zetamatta/go-texts/mbcs"
	"github.com/zetamatta/nyagos/dos"
)

func CmdMsgBox(args []any_t) []any_t {
	var message string
	title := "nyagos"
	if len(args) >= 1 {
		message = fmt.Sprint(args[0])
	}
	if len(args) >= 2 {
		title = fmt.Sprint(args[1])
	}
	msgbox.Show(0, message, title, msgbox.OK)
	return []any_t{}
}

func CmdElevated([]any_t) []any_t {
	flag, _ := dos.IsElevated()
	return []any_t{flag}
}

func CmdNetDriveToUNC(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{}
	}
	path, ok := args[0].(string)
	if !ok {
		return []any_t{path}
	}
	unc := dos.NetDriveToUNC(path)
	return []any_t{unc}
}

func CmdShellExecute(args []any_t) []any_t {
	err := dos.ShellExecute(
		toStr(args, 0),
		toStr(args, 1),
		toStr(args, 2),
		toStr(args, 3))
	if err != nil {
		return []any_t{nil, err}
	} else {
		return []any_t{true}
	}
}

func CmdRawExec(this *Param) []any_t {
	argv := stackToSlice(this)
	xcmd := exec.Command(argv[0], argv[1:]...)
	xcmd.Stdin = this.In
	xcmd.Stdout = this.Out
	xcmd.Stderr = this.Err
	err := xcmd.Run()
	errorlevel, errorlevelOk := dos.GetErrorLevel(xcmd)
	if !errorlevelOk {
		errorlevel = 255
	}
	if err != nil {
		fmt.Fprintln(xcmd.Stderr, err.Error())
		return []any_t{errorlevel, err.Error()}
	} else {
		return []any_t{errorlevel}
	}
}

func CmdAtoU(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, TooFewArguments}
	}
	if s, ok := args[0].(string); ok {
		if val, err := mbcs.AtoU([]byte(s), mbcs.ConsoleCP()); err == nil {
			return []any_t{val}
		} else {
			return []any_t{nil, err}
		}
	} else {
		return []any_t{fmt.Sprint(args[0])}
	}
}

func CmdUtoA(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, TooFewArguments}
	}
	utf8 := fmt.Sprint(args[0])
	bin, err := mbcs.UtoA(utf8, mbcs.ConsoleCP(), true)
	if err != nil {
		return []any_t{nil, err}
	}
	return []any_t{bin, nil}
}
