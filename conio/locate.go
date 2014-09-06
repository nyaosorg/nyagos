package conio

/*
#include <windows.h>

void getLocate(HANDLE hConin,short* X,short* Y){
	CONSOLE_SCREEN_BUFFER_INFO csbi;
	GetConsoleScreenBufferInfo(hConin,&csbi);
	*X = (int)csbi.dwCursorPosition.X;
	*Y = (int)csbi.dwCursorPosition.Y;
}
*/
import "C"

var hConout = C.GetStdHandle(C.STD_OUTPUT_HANDLE)

func GetLocate() (int, int) {
	var x C.short
	var y C.short
	C.getLocate(hConout, &x, &y)
	return int(x), int(y)
}

func Locate(x, y int) {
	C.SetConsoleCursorPosition(hConout, C.COORD{X: C.SHORT(x), Y: C.SHORT(y)})
}
