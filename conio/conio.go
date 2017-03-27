package conio

import (
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32")

var hConout, _ = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)

type coord_t struct {
	X int16
	Y int16
}

type small_rect_t struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

type console_screen_buffer_info_t struct {
	Size              coord_t
	CursorPosition    coord_t
	Attributes        uint16
	Window            small_rect_t
	MaximumWindowSize coord_t
}

var getConsoleScreenBufferInfo = kernel32.NewProc("GetConsoleScreenBufferInfo")

func GetScreenBufferInfo() *console_screen_buffer_info_t {
	var csbi console_screen_buffer_info_t
	getConsoleScreenBufferInfo.Call(
		uintptr(hConout),
		uintptr(unsafe.Pointer(&csbi)))
	return &csbi
}

func (this *console_screen_buffer_info_t) ViewSize() (int, int) {
	return int(this.Window.Right-this.Window.Left) + 1,
		int(this.Window.Bottom-this.Window.Top) + 1
}

func (this *console_screen_buffer_info_t) CursorPos() (int, int) {
	return int(this.CursorPosition.X), int(this.CursorPosition.Y)
}

func GetLocate() (int, int) {
	return GetScreenBufferInfo().CursorPos()
}
