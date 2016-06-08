package getch

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32")

const (
	RIGHT_ALT_PRESSED  = 1
	LEFT_ALT_PRESSED   = 2
	RIGHT_CTRL_PRESSED = 4
	LEFT_CTRL_PRESSED  = 8
	CTRL_PRESSED       = RIGHT_CTRL_PRESSED | LEFT_CTRL_PRESSED
	ALT_PRESSED        = RIGHT_ALT_PRESSED | LEFT_ALT_PRESSED
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
	KeyCode    rune
	ScanCode   uint16
	ShiftState uint32
}

func ctrlCHandler(ch chan os.Signal) {
	for _ = range ch {
		keyBuffer = append(keyBuffer, keyInfo{3, 0, LEFT_CTRL_PRESSED})
	}
}

func DisableCtrlC() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go ctrlCHandler(ch)
}

func getKeys() []keyInfo {
	var numberOfEventsRead uint32
	var events [10]inputRecordT
	var orgConMode uint32

	result := make([]keyInfo, 0, 0)

	getConsoleMode.Call(uintptr(hConin),
		uintptr(unsafe.Pointer(&orgConMode)))
	setConsoleMode.Call(uintptr(hConin), 0)
	var precode rune = 0
	for len(result) <= 0 {
		readConsoleInput.Call(
			uintptr(hConin),
			uintptr(unsafe.Pointer(&events[0])),
			uintptr(len(events)),
			uintptr(unsafe.Pointer(&numberOfEventsRead)))
		for i := uint32(0); i < numberOfEventsRead; i++ {
			if events[i].eventType == KEY_EVENT && events[i].bKeyDown != 0 {
				var keycode = rune(events[i].unicodeChar)
				if keycode != 0 {
					if precode != 0 {
						keycode = utf16.DecodeRune(precode, keycode)
						precode = 0
					} else if utf16.IsSurrogate(keycode) {
						precode = keycode
						continue
					}
				}
				result = append(result, keyInfo{
					keycode,
					events[i].wVirtualKeyCode,
					events[i].dwControlKeyState,
				})
			}
		}
	}
	setConsoleMode.Call(uintptr(hConin), uintptr(orgConMode))
	return result
}

var keyBuffer []keyInfo
var keyBufferRead = 0

func getKey() (rune, uint16, uint32) {
	for keyBuffer == nil || keyBufferRead >= len(keyBuffer) {
		keyBuffer = getKeys()
		keyBufferRead = 0
	}
	r := keyBuffer[keyBufferRead]
	keyBufferRead++
	return r.KeyCode, r.ScanCode, r.ShiftState
}

func Full() (rune, uint16, uint32) {
	code, scan, shift := getKey()
	if code < 0xDC00 || 0xDFFF < code {
		return code, scan, shift
	}
	code2, _, _ := getKey()
	return utf16.DecodeRune(code, code2), scan, shift
}

func Rune() rune {
	for {
		ch, _, _ := Full()
		if ch != 0 {
			return ch
		}
	}
}
