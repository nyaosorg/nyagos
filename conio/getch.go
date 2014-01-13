package conio

/*
#include <windows.h>
KEY_EVENT_RECORD* getKeyEvent(struct _INPUT_RECORD *input) {
	return &(input->Event.KeyEvent);
}
char getAsciiChar(struct _INPUT_RECORD *input) {
	return input->Event.KeyEvent.uChar.AsciiChar;
}
WCHAR getUnicodeChar(struct _INPUT_RECORD *input) {
	return input->Event.KeyEvent.uChar.UnicodeChar;
}
WCHAR getVirtualKeyCode(struct _INPUT_RECORD *input) {
	return input->Event.KeyEvent.wVirtualKeyCode;
}
*/
import "C"
import "unicode/utf16"

var hConin C.HANDLE = C.GetStdHandle(C.STD_INPUT_HANDLE)

var buffer [10][2]uint16
var readptr int = 0
var stacked int = 0

func GetKey() (rune, uint16) {
	if readptr >= stacked {
		readptr = 0
		stacked = 0
		for stacked <= 0 {
			var numberOfEventsRead C.DWORD
			var events [len(buffer)]C.struct__INPUT_RECORD
			var orgConMode C.DWORD

			C.GetConsoleMode(hConin, &orgConMode)
			C.SetConsoleMode(hConin, 0)
			C.ReadConsoleInputW(hConin,
				&events[0],
				C.DWORD(len(events)),
				&numberOfEventsRead)
			C.SetConsoleMode(hConin, orgConMode)
			for i := C.DWORD(0); i < numberOfEventsRead; i++ {
				if events[i].EventType == C.KEY_EVENT && C.getKeyEvent(&events[i]).bKeyDown != C.FALSE {
					buffer[stacked][0] = uint16(C.getUnicodeChar(&events[i]))
					buffer[stacked][1] = uint16(C.getVirtualKeyCode(&events[i]))
					stacked++
				}
			}
		}
	}
	rc := buffer[readptr]
	readptr++
	if rc[0] == 0 {
		return rune(0), rc[1]
	} else {
		return (utf16.Decode([]uint16{rc[0]}))[0], rc[1]
	}
}

func GetCh() rune {
	for {
		ch, _ := GetKey()
		if ch != 0 {
			return ch
		}
	}
}
