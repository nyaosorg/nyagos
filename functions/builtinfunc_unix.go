// +build !windows

package functions

func CmdMsgBox(args []any_t) []any_t {
	return []any_t{nil, "not support"}
}

func CmdElevated([]any_t) []any_t {
	return []any_t{nil, "not support"}
}
func CmdShellExecute(args []any_t) []any_t {
	return []any_t{nil, "not support"}
}
func CmdRawExec(this *Param) []any_t {
	return []any_t{nil, "not support"}
}
func CmdAtoU(args []any_t) []any_t {
	return []any_t{nil, "not support"}
}
func CmdUtoA(args []any_t) []any_t {
	return []any_t{nil, "not support"}
}
