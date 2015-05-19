package dos

import "syscall"

const (
	COINIT_APARTMENTTHREADED uintptr = 0x2
	COINIT_MULTITHREADED     uintptr = 0x0
	COINIT_DISABLE_OLE1DDE   uintptr = 0x4
	COINIT_SPEED_OVER_MEMORY uintptr = 0x8
)

var ole32 = syscall.NewLazyDLL("ole32")
var coInitializeEx = ole32.NewProc("CoInitializeEx")
var coUninitialize = ole32.NewProc("CoUninitialize")

func CoInitializeEx(res, opt uintptr) {
	coInitializeEx.Call(res, opt)
}

func CoUninitialize() {
	coUninitialize.Call()
}
