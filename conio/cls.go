package conio

/*
#include <windows.h>

DWORD getSize(CONSOLE_SCREEN_BUFFER_INFO *csbi){
	return csbi->dwSize.X * csbi->dwSize.Y;
}

DWORD getWidth(CONSOLE_SCREEN_BUFFER_INFO *csbi){
	return csbi->srWindow.Right - csbi->srWindow.Left;
}

DWORD getHeight(CONSOLE_SCREEN_BUFFER_INFO *csbi){
	return csbi->srWindow.Bottom - csbi->srWindow.Top;
}
*/
import "C"

func GetScreenSize() (int, int) {
	var csbi C.CONSOLE_SCREEN_BUFFER_INFO
	C.GetConsoleScreenBufferInfo(hConout, &csbi)
	return int(C.getWidth(&csbi)), int(C.getHeight(&csbi))
}

func Cls() {
	var csbi C.CONSOLE_SCREEN_BUFFER_INFO
	var coordScreen C.COORD
	var cCharsWritten C.DWORD

	C.GetConsoleScreenBufferInfo(hConout, &csbi)
	dwConSize := C.getSize(&csbi)

	coordScreen.X = 0
	coordScreen.Y = 0
	C.FillConsoleOutputCharacter(hConout,
		C.CHAR(' '),
		dwConSize,
		coordScreen,
		&cCharsWritten)

	C.FillConsoleOutputAttribute(hConout,
		csbi.wAttributes,
		dwConSize,
		coordScreen,
		&cCharsWritten)

	C.SetConsoleCursorPosition(hConout, coordScreen)
}
