package conio

import (
	"syscall"
	"unsafe"
)

type coord_t struct {
	X int16
	Y int16
}

func (this *coord_t) Pack() uintptr {
	return *(*uintptr)(unsafe.Pointer(this))
}

type small_rect_t struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

type consoleScreenBufferInfo struct {
	Size              coord_t
	CursorPosition    coord_t
	Attributes        uint16
	Window            small_rect_t
	MaximumWindowSize coord_t
}

func (this *consoleScreenBufferInfo) Address() uintptr {
	return uintptr(unsafe.Pointer(this))
}

var hConout, _ = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
var getConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
var setConsoleCursorPosition = kernel32.NewProc("SetConsoleCursorPosition")

func GetLocate() (int, int) {
	var csbi consoleScreenBufferInfo
	getConsoleScreenBufferInfo.Call(uintptr(hConout), csbi.Address())
	return int(csbi.CursorPosition.X), int(csbi.CursorPosition.Y)
}

func Locate(x, y int) {
	csbi := coord_t{X: int16(x), Y: int16(y)}
	setConsoleCursorPosition.Call(uintptr(hConout), csbi.Pack())
}
