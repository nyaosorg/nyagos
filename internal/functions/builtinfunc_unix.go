//go:build !windows
// +build !windows

package functions

func CmdMsgBox(args []anyT) []anyT {
	return []anyT{nil, "not support"}
}

func CmdElevated([]anyT) []anyT {
	return []anyT{nil, "not support"}
}
func CmdShellExecute(args []anyT) []anyT {
	return []anyT{nil, "not support"}
}
func CmdRawExec(this *Param) []anyT {
	return []anyT{nil, "not support"}
}
func CmdAtoU(args []anyT) []anyT {
	return []anyT{nil, "not support"}
}
func CmdUtoA(args []anyT) []anyT {
	return []anyT{nil, "not support"}
}
