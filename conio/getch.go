package conio

import "syscall"
import "unsafe"
import "unicode/utf16"

type inputRecordT struct {
	eventType uint16
	_         uint16
	// _KEY_EVENT_RECORD {
	bKeyDown         uintptr
	wRepeartCount    uint16
	wVirtualKeyCode  uint16
	wVirtualScanCode uint16
	unicodeChar      uint16
	// }
	dwControlKeyState uint32
}

var buffer [10][2]uint16
var readptr = 0
var stacked = 0

var createFile = kernel32.NewProc("CreateFileW")
var getConsoleMode = kernel32.NewProc("GetConsoleMode")
var setConsoleMode = kernel32.NewProc("SetConsoleMode")
var readConsoleInput = kernel32.NewProc("ReadConsoleInputW")

var conioS, _ = syscall.UTF16FromString("CONIN$")
var hConin, _, _ = createFile.Call(
	uintptr(unsafe.Pointer(&conioS[0])),
	GENERIC_READ,
	FILE_SHARE_READ,
	0,
	OPEN_EXISTING,
	FILE_ATTRIBUTE_NORMAL,
	0)

func GetKey() (rune, uint16) {
	if readptr >= stacked {
		readptr = 0
		stacked = 0
		for stacked <= 0 {
			var numberOfEventsRead uint32
			var events [len(buffer)]inputRecordT
			var orgConMode uint32

			getConsoleMode.Call(uintptr(hConin), uintptr(unsafe.Pointer(&orgConMode)))
			setConsoleMode.Call(uintptr(hConin), 0)
			readConsoleInput.Call(
				uintptr(hConin),
				uintptr(unsafe.Pointer(&events[0])),
				uintptr(len(events)),
				uintptr(unsafe.Pointer(&numberOfEventsRead)))
			setConsoleMode.Call(uintptr(hConin), uintptr(orgConMode))
			for i := uint32(0); i < numberOfEventsRead; i++ {
				if events[i].eventType == KEY_EVENT && events[i].bKeyDown != 0 {
					buffer[stacked][0] = events[i].unicodeChar
					buffer[stacked][1] = events[i].wVirtualKeyCode
					stacked++
				} else {
					if CtrlC {
						buffer[stacked][0] = 3
						buffer[stacked][1] = 0
						stacked++
						CtrlC = false
					}
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
