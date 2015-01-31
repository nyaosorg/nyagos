package conio

import (
	"syscall"
)

var hConout, _ = syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)

func GetLocate() (int, int) {
	csbi := GetScreenBufferInfo()
	return int(csbi.CursorPosition.X), int(csbi.CursorPosition.Y)
}

func Locate(x, y int) {
	csbi := coord_t{X: int16(x), Y: int16(y)}
	setConsoleCursorPosition.Call(uintptr(hConout), csbi.Pack())
}
