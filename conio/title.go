package conio

//#include <windows.h>
import "C"
import "syscall"

func SetTitle(title string) {
	ctitle, err := syscall.UTF16FromString(title)
	if ctitle != nil && err == nil {
		C.SetConsoleTitleW((*C.WCHAR)(&ctitle[0]))
	}
}
