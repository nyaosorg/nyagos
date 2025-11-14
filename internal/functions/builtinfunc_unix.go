//go:build !windows
// +build !windows

package functions

func (*Env) CmdMsgBox(args []any) []any {
	return []any{nil, "not support"}
}

func (*Env) CmdElevated([]any) []any {
	return []any{nil, "not support"}
}
func (*Env) CmdShellExecute(args []any) []any {
	return []any{nil, "not support"}
}
func (*Env) CmdRawExec(this *Param) []any {
	return []any{nil, "not support"}
}
func (*Env) CmdAtoU(args []any) []any {
	return args
}
func (*Env) CmdUtoA(args []any) []any {
	return args
}
func (*Env) CmdAnsiToUtf8IfNeeded(args []any) []any {
	return args
}
