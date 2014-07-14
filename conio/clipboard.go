package conio

/*
#include "windows.h"

int readClipBoard( wchar_t *buffer , size_t size )
{
	int len=0;
	OpenClipboard(NULL);
    HANDLE hText=GetClipboardData(CF_UNICODETEXT);
    if( hText != NULL ){
        wchar_t *pText = (wchar_t*)GlobalLock(hText);
        if( pText != NULL ){
			while( pText[len] != 0 )
				len++;
			if( len >= size-1 ){
				len = size-1;
				pText[len] = 0;
			}
			memcpy( buffer,pText,(len+1)*sizeof(wchar_t));
            GlobalUnlock(hText);
        }
    }
    CloseClipboard();
	return len;
}
*/
import "C"
import "unsafe"
import "syscall"

func GetClipboardUtf16() []uint16 {
	var text [syscall.MAX_PATH]uint16
	length := C.readClipBoard((*C.wchar_t)(unsafe.Pointer(&text[0])), (C.size_t)(len(text)))
	return text[:length]
}

func GetClipboardString() string {
	return syscall.UTF16ToString(GetClipboardUtf16())
}
