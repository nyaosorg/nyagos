//go:build !windows
// +build !windows

package functions

func CmdMsgBox(args []any) []any {
	return []any{nil, "not support"}
}

func CmdElevated([]any) []any {
	return []any{nil, "not support"}
}
func CmdShellExecute(args []any) []any {
	return []any{nil, "not support"}
}
func CmdRawExec(this *Param) []any {
	return []any{nil, "not support"}
}
func CmdAtoU(args []any) []any {
	return args
}
func CmdUtoA(args []any) []any {
	return args
}
func CmdAnsiToUtf8IfNeeded(args []any) []any {
	return args
}
