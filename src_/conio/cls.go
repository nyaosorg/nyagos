package conio

import "unsafe"

var fillConsoleOutputCharacter = kernel32.NewProc("FillConsoleOutputCharacterW")
var fillConsoleOutputAttribute = kernel32.NewProc("FillConsoleOutputAttribute")

func Cls() {
	csbi := GetScreenBufferInfo()
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
