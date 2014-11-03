package conio

import "syscall"

var setConsoleCtrlHandler = kernel32.NewProc("SetConsoleCtrlHandler")

var list = []func(){}

func callBack(dwCtrlType uintptr) uintptr {
	switch dwCtrlType {
	case CTRL_CLOSE_EVENT, CTRL_LOGOFF_EVENT, CTRL_SHUTDOWN_EVENT:
		for i := len(list) - 1; i >= 0; i-- {
			list[i]()
		}
	}
	return 0
}

func OnClose(f func()) {
	if len(list) <= 0 {
		setConsoleCtrlHandler.Call(
			/* syscall.NewCallback(callBack),*/
			syscall.NewCallbackCDecl(callBack),
			uintptr(1))
	}
	list = append(list, f)
}
