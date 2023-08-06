package functions

import (
	"fmt"
	"os/exec"
	"unicode/utf8"

	"github.com/mattn/msgbox"
	"github.com/nyaosorg/go-windows-mbcs"
	"github.com/nyaosorg/go-windows-su"
)

func CmdMsgBox(args []any) []any {
	var message string
	title := "nyagos"
	if len(args) >= 1 {
		message = fmt.Sprint(args[0])
	}
	if len(args) >= 2 {
		title = fmt.Sprint(args[1])
	}
	msgbox.Show(0, message, title, msgbox.OK)
	return []any{}
}

func CmdElevated([]any) []any {
	flag, _ := su.IsElevated()
	return []any{flag}
}

func CmdShellExecute(args []any) []any {
	pid, err := su.ShellExecute(
		toStr(args, 0),
		toStr(args, 1),
		toStr(args, 2),
		toStr(args, 3))
	if err != nil {
		return []any{nil, err}
	}
	return []any{pid}
}

func CmdRawExec(this *Param) []any {
	argv := stackToSlice(this)
	xcmd := exec.Command(argv[0], argv[1:]...)
	xcmd.Stdin = this.In
	xcmd.Stdout = this.Out
	xcmd.Stderr = this.Err
	err := xcmd.Run()
	errorlevel := xcmd.ProcessState.ExitCode()
	if err != nil {
		fmt.Fprintln(xcmd.Stderr, err.Error())
		return []any{errorlevel, err.Error()}
	}
	return []any{errorlevel}
}

func CmdAtoU(args []any) []any {
	if len(args) < 1 {
		return []any{nil, TooFewArguments}
	}
	if s, ok := args[0].(string); ok {
		val, err := mbcs.AnsiToUtf8([]byte(s), mbcs.ConsoleCP())
		if err != nil {
			return []any{nil, err}
		}
		return []any{val}
	}
	return []any{fmt.Sprint(args[0])}
}

func CmdUtoA(args []any) []any {
	if len(args) < 1 {
		return []any{nil, TooFewArguments}
	}
	utf8 := fmt.Sprint(args[0])
	bin, err := mbcs.Utf8ToAnsi(utf8, mbcs.ConsoleCP())
	if err != nil {
		return []any{nil, err}
	}
	return []any{bin, nil}
}

func CmdAnsiToUtf8IfNeeded(args []any) []any {
	if len(args) < 1 {
		return []any{nil, TooFewArguments}
	}
	if s, ok := args[0].(string); ok {
		if utf8.ValidString(s) {
			return []any{s}
		}
		val, err := mbcs.AnsiToUtf8([]byte(s), mbcs.ConsoleCP())
		if err != nil {
			return []any{nil, err}
		}
		return []any{val}
	}
	return []any{fmt.Sprint(args[0])}
}
