package nodos

import (
	"golang.org/x/sys/windows"
)

var ole32 = windows.NewLazySystemDLL("ole32.dll")
var procCoInitializeEx = ole32.NewProc("CoInitializeEx")
var procCoUninitialize = ole32.NewProc("CoUninitialize")

func coInitializeEx(res uintptr, opt uintptr) {
	procCoInitializeEx.Call(res, opt)
}

func coUninitialize() {
	procCoUninitialize.Call()
}
