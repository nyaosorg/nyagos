package conio

/*
#include <windows.h>

extern void callBack();
static WINAPI myHandleRoutine( DWORD dwCtrlType )
{
    switch( dwCtrlType ){
    case CTRL_CLOSE_EVENT:
    case CTRL_LOGOFF_EVENT:
    case CTRL_SHUTDOWN_EVENT:
		callBack();
        break;
    default:
        break;
    }
    return FALSE;
}

static void MySetConsoleCtrlHandler()
{
	SetConsoleCtrlHandler( myHandleRoutine , TRUE );
}

*/
import "C"

var list = []func(){}

//export callBack
func callBack() {
	for _, f := range list {
		f()
	}
}

func OnClose(f func()) {
	if len(list) <= 0 {
		C.MySetConsoleCtrlHandler()
	}
	list = append(list, f)
}
