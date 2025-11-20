//go:build !windows
// +build !windows

package functions

func CmdMsgBox(param *Param) []any {
	return []any{nil, "not support"}
}

func CmdElevated(param *Param) []any {
	return []any{nil, "not support"}
}
func CmdShellExecute(param *Param) []any {
	return []any{nil, "not support"}
}
func CmdRawExec(param *Param) []any {
	return []any{nil, "not support"}
}
func CmdAtoU(param *Param) []any {
	return param.Args
}
func CmdUtoA(param *Param) []any {
	return param.Args
}
func CmdAnsiToUtf8IfNeeded(param *Param) []any {
	return param.Args
}
