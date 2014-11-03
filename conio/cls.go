package conio

import "unsafe"

func GetScreenSize() (int, int) {
	var csbi consoleScreenBufferInfo
	getConsoleScreenBufferInfo.Call(uintptr(hConout), csbi.Address())
	return int(csbi.Size.X), int(csbi.Size.Y)
}

var fillConsoleOutputCharacter = kernel32.NewProc("FillConsoleOutputCharacterW")
var fillConsoleOutputAttribute = kernel32.NewProc("FillConsoleOutputAttribute")

func Cls() {
	var csbi consoleScreenBufferInfo
	getConsoleScreenBufferInfo.Call(uintptr(hConout), csbi.Address())
	var cCharsWritten uint32
	c := coord_t{0, 0}
	coordScreen := c.Pack()
	dwConSize := csbi.Size.Pack()

	fillConsoleOutputCharacter.Call(
		uintptr(hConout),
		uintptr(' '),
		dwConSize,
		coordScreen,
		uintptr(unsafe.Pointer(&cCharsWritten)))

	fillConsoleOutputAttribute.Call(
		uintptr(hConout),
		uintptr(csbi.Attributes),
		dwConSize,
		coordScreen,
		uintptr(unsafe.Pointer(&cCharsWritten)))

	setConsoleCursorPosition.Call(uintptr(hConout), coordScreen)
}
