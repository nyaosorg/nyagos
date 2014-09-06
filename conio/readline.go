package conio

type KeyFuncResult int

const (
	CONTINUE KeyFuncResult = iota
	ENTER    KeyFuncResult = iota
	ABORT    KeyFuncResult = iota
)

type KeyFuncT interface {
	Call(buffer *ReadLineBuffer) KeyFuncResult
}

type KeyGoFuncT struct {
	F func(buffer *ReadLineBuffer) KeyFuncResult
}

func (this *KeyGoFuncT) Call(buffer *ReadLineBuffer) KeyFuncResult {
	return this.F(buffer)
}

var KeyMap = map[rune]KeyFuncT{
	'\b':         &KeyGoFuncT{KeyFuncBackSpace},
	'\r':         &KeyGoFuncT{KeyFuncEnter},
	'\x7F':       &KeyGoFuncT{KeyFuncDelete},
	('A' & 0x1f): &KeyGoFuncT{KeyFuncHead},
	('B' & 0x1f): &KeyGoFuncT{KeyFuncBackword},
	('D' & 0x1f): &KeyGoFuncT{KeyFuncDeleteOrAbort},
	('E' & 0x1f): &KeyGoFuncT{KeyFuncTail},
	('F' & 0x1f): &KeyGoFuncT{KeyFuncForward},
	('K' & 0x1f): &KeyGoFuncT{KeyFuncClearAfter},
	('L' & 0x1F): &KeyGoFuncT{KeyFuncCLS},
	('U' & 0x1F): &KeyGoFuncT{KeyFuncClearBefore},
	('Y' & 0x1F): &KeyGoFuncT{KeyFuncPaste},
	('[' & 0x1F): &KeyGoFuncT{KeyFuncClear},
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

var ZeroMap = map[uint16]KeyFuncT{
	K_LEFT:  &KeyGoFuncT{KeyFuncBackword},
	K_RIGHT: &KeyGoFuncT{KeyFuncForward},
	K_DEL:   &KeyGoFuncT{KeyFuncDelete},
	K_HOME:  &KeyGoFuncT{KeyFuncHead},
	K_END:   &KeyGoFuncT{KeyFuncTail},
	K_CTRL:  &KeyGoFuncT{KeyFuncPass},
	K_SHIFT: &KeyGoFuncT{KeyFuncPass},
}

func ReadLine(prompt_ func() int) (string, KeyFuncResult) {
	this := ReadLineBuffer{Buffer: make([]rune, 20)}
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
		var f KeyFuncT
		var ok bool
		if this.Unicode != 0 {
			f, ok = KeyMap[this.Unicode]
			if !ok {
				//f = KeyFuncInsertReport
				f = &KeyGoFuncT{KeyFuncInsertSelf}
			}
		} else {
			f, ok = ZeroMap[this.Keycode]
			if !ok {
				f = &KeyGoFuncT{KeyFuncPass}
			}
		}
		rc := f.Call(&this)
		if rc != CONTINUE {
			stdOut.WriteRune('\n')
			stdOut.Flush()
			return this.String(), rc
		}
	}
}
