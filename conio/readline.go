package conio

type KeyFuncResult int

const (
	CONTINUE KeyFuncResult = iota
	ENTER    KeyFuncResult = iota
	ABORT    KeyFuncResult = iota
)

var KeyMap = map[rune]func(*ReadLineBuffer) KeyFuncResult{
	'\b':         KeyFuncBackSpace,
	'\r':         KeyFuncEnter,
	'\x7F':       KeyFuncDelete,
	('A' & 0x1f): KeyFuncHead,
	('B' & 0x1f): KeyFuncBackword,
	('D' & 0x1f): KeyFuncDeleteOrAbort,
	('E' & 0x1f): KeyFuncTail,
	('F' & 0x1f): KeyFuncForward,
	('K' & 0x1f): KeyFuncClearAfter,
	('L' & 0x1F): KeyFuncCLS,
	('U' & 0x1F): KeyFuncClearBefore,
	('Y' & 0x1F): KeyFuncPaste,
	('[' & 0x1F): KeyFuncClear,
}

// KeyCode from
// http://msdn.microsoft.com/ja-jp/library/windows/desktop/dd375731(v=vs.85).aspx
const (
	K_LEFT  = 0x25
	K_RIGHT = 0x27
	K_DEL   = 0x2E
	K_HOME  = 0x24
	K_END   = 0x23
	K_CTRL  = 0x11
	K_SHIFT = 0x10
	K_UP    = 0x26
	K_DOWN  = 0x28
)

var ZeroMap = map[uint16]func(*ReadLineBuffer) KeyFuncResult{
	K_LEFT:  KeyFuncBackword,
	K_RIGHT: KeyFuncForward,
	K_DEL:   KeyFuncDelete,
	K_HOME:  KeyFuncHead,
	K_END:   KeyFuncTail,
	K_CTRL:  KeyFuncPass,
	K_SHIFT: KeyFuncPass,
}

func ReadLine(prompt_ func() int) (string, KeyFuncResult) {
	var this ReadLineBuffer
	this.Buffer = make([]rune, 20)
	this.Length = 0
	this.Cursor = 0
	this.ViewStart = 0
	this.ViewWidth, _ = GetScreenSize()
	this.ViewWidth--

	this.Prompt = prompt_
	if this.Prompt != nil {
		this.ViewWidth = this.ViewWidth - this.Prompt()
	}
	for {
		stdOut.Flush()
		shineCursor()
		this.Unicode, this.Keycode = GetKey()
		var f func(*ReadLineBuffer) KeyFuncResult
		var ok bool
		if this.Unicode != 0 {
			f, ok = KeyMap[this.Unicode]
			if !ok {
				//f = KeyFuncInsertReport
				f = KeyFuncInsertSelf
			}
		} else {
			f, ok = ZeroMap[this.Keycode]
			if !ok {
				f = KeyFuncPass
			}
		}
		rc := f(&this)
		if rc != CONTINUE {
			stdOut.WriteRune('\n')
			stdOut.Flush()
			return this.String(), rc
		}
	}
}
