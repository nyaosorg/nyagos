package conio

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type inputRecordT struct {
	eventType uint16
	_         uint16
	// _KEY_EVENT_RECORD {
	bKeyDown         int32
	wRepeartCount    uint16
	wVirtualKeyCode  uint16
	wVirtualScanCode uint16
	unicodeChar      uint16
	// }
	dwControlKeyState uint32
}

var getConsoleMode = kernel32.NewProc("GetConsoleMode")
var setConsoleMode = kernel32.NewProc("SetConsoleMode")
var readConsoleInput = kernel32.NewProc("ReadConsoleInputW")

var hConin syscall.Handle

func init() {
	var err error
	hConin, err = syscall.Open("CONIN$", syscall.O_RDWR, 0)
	if err != nil {
		panic(fmt.Sprintf("conio: %v", err))
	}
}

type keyInfo struct {
	KeyCode  rune
	ScanCode uint16
}

func DisableCtrlC() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		for _ = range ch {
			if keyPipe != nil {
				go func() {
					keyPipe <- keyInfo{3, 0}
				}()
			}
		}
	}()
}

var keyPipe chan keyInfo = nil

func keyGoRuntine(pipe chan keyInfo) {
	var numberOfEventsRead uint32
	var events [10]inputRecordT
	var orgConMode uint32

	getConsoleMode.Call(uintptr(hConin),
		uintptr(unsafe.Pointer(&orgConMode)))
	setConsoleMode.Call(uintptr(hConin), uintptr(ENABLE_PROCESSED_INPUT))
	readConsoleInput.Call(
		uintptr(hConin),
		uintptr(unsafe.Pointer(&events[0])),
		uintptr(len(events)),
		uintptr(unsafe.Pointer(&numberOfEventsRead)))
	setConsoleMode.Call(uintptr(hConin), uintptr(orgConMode))
	for i := uint32(0); i < numberOfEventsRead; i++ {
		if events[i].eventType == KEY_EVENT && events[i].bKeyDown != 0 {
			var keycode rune
			if events[i].unicodeChar == 0 {
				keycode = rune(0)
			} else {
				keycode = utf16.Decode([]uint16{events[i].unicodeChar})[0]
			}
			pipe <- keyInfo{
				keycode,
				events[i].wVirtualKeyCode,
			}
		}
	}
	// Not to read keyboard data on not requested time
	// (ex. other application is running)
	// shutdown goroutine.
	pipe <- keyInfo{0, 0}
}

func GetKey() (rune, uint16) {
	if keyPipe == nil {
		keyPipe = make(chan keyInfo, 10)
		go keyGoRuntine(keyPipe)
	}
	for {
		keyInfo := <-keyPipe
		if keyInfo.KeyCode != 0 || keyInfo.ScanCode != 0 {
			return keyInfo.KeyCode, keyInfo.ScanCode
		}
		// When keyGoRuntine has shutdowned, restart.
		go keyGoRuntine(keyPipe)
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
