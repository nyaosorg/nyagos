package conio

import "unsafe"

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

var getConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")
var setConsoleCursorPosition = kernel32.NewProc("SetConsoleCursorPosition")

func GetScreenBufferInfo() *consoleScreenBufferInfo {
	var csbi consoleScreenBufferInfo
	getConsoleScreenBufferInfo.Call(
		uintptr(hConout),
		uintptr(unsafe.Pointer(&csbi)))
	return &csbi
}

func (this *consoleScreenBufferInfo) ViewSize() (int, int) {
	return int(this.Window.Right-this.Window.Left) + 1,
		int(this.Window.Bottom-this.Window.Top) + 1
}
