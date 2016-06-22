package dos

import "syscall"

var ole32 = syscall.NewLazyDLL("ole32")
var coInitializeEx = ole32.NewProc("CoInitializeEx")
var coUninitialize = ole32.NewProc("CoUninitialize")

func CoInitializeEx(res, opt uintptr) {
	coInitializeEx.Call(res, opt)
}

func CoUninitialize() {
	coUninitialize.Call()
}
