package readline

import (
	"fmt"
	"strings"

	"../../conio"
)

type Result int

const (
	CONTINUE Result = iota
	ENTER    Result = iota
	ABORT    Result = iota
)

type KeyFuncT interface {
	Call(buffer *Buffer) Result
}

type KeyGoFuncT struct {
	F func(buffer *Buffer) Result
}

func (this *KeyGoFuncT) Call(buffer *Buffer) Result {
	return this.F(buffer)
}

var keyMap = map[rune]KeyFuncT{
	name2char[K_CTRL_A]: name2func[F_BEGINNING_OF_LINE],
	name2char[K_CTRL_B]: name2func[F_BACKWARD_CHAR],
	name2char[K_CTRL_C]: name2func[F_INTR],
	name2char[K_CTRL_D]: name2func[F_DELETE_OR_ABORT],
	name2char[K_CTRL_E]: name2func[F_END_OF_LINE],
	name2char[K_CTRL_F]: name2func[F_FORARD_CHAR],
	name2char[K_CTRL_H]: name2func[F_BACKWARD_DELETE_CHAR],
	name2char[K_CTRL_K]: name2func[F_KILL_LINE],
	name2char[K_CTRL_L]: name2func[F_CLEAR_SCREEN],
	name2char[K_CTRL_M]: name2func[F_ACCEPT_LINE],
	name2char[K_CTRL_U]: name2func[F_UNIX_LINE_DISCARD],
	name2char[K_CTRL_Y]: name2func[F_YANK],
	name2char[K_DELETE]: name2func[F_DELETE_CHAR],
	name2char[K_ENTER]:  name2func[F_ACCEPT_LINE],
	name2char[K_ESCAPE]: name2func[F_KILL_WHOLE_LINE],
}

var scanMap = map[uint16]KeyFuncT{
	name2scan[K_CTRL]:   name2func[F_PASS],
	name2scan[K_DELETE]: name2func[F_DELETE_CHAR],
	name2scan[K_END]:    name2func[F_END_OF_LINE],
	name2scan[K_HOME]:   name2func[F_BEGINNING_OF_LINE],
	name2scan[K_LEFT]:   name2func[F_BACKWARD_CHAR],
	name2scan[K_RIGHT]:  name2func[F_FORARD_CHAR],
	name2scan[K_SHIFT]:  name2func[F_PASS],
}

var altMap = map[uint16]KeyFuncT{
	name2alt[K_ALT_V]: name2func[F_YANK],
}

func normWord(src string) string {
	return strings.Replace(strings.ToUpper(src), "-", "_", -1)
}

func BindKeyFunc(keyName string, funcValue KeyFuncT) error {
	keyName_ := normWord(keyName)
	if altValue, altOk := name2alt[keyName_]; altOk {
		altMap[altValue] = funcValue
		return nil
	} else if charValue, charOk := name2char[keyName_]; charOk {
		keyMap[charValue] = funcValue
		return nil
	} else if scanValue, scanOk := name2scan[keyName_]; scanOk {
		scanMap[scanValue] = funcValue
		return nil
	} else {
		return fmt.Errorf("%s: no such keyname", keyName)
	}
}

func GetFunc(funcName string) (KeyFuncT, error) {
	rc, ok := name2func[normWord(funcName)]
	if ok {
		return rc, nil
	} else {
		return nil, fmt.Errorf("%s: not found in the function-list", funcName)
	}
}

func BindKeySymbol(keyName, funcName string) error {
	funcValue, funcOk := name2func[normWord(funcName)]
	if !funcOk {
		return fmt.Errorf("%s: no such function.", funcName)
	}
	return BindKeyFunc(keyName, funcValue)
}

func BindKeySymbolFunc(keyName, funcName string, funcValue KeyFuncT) error {
	name2func[normWord(funcName)] = funcValue
	return BindKeyFunc(keyName, funcValue)
}

func ReadLine(prompt_ func() int) (string, Result) {
	this := Buffer{Buffer: make([]rune, 20)}
	this.ViewWidth, _ = conio.GetScreenBufferInfo().ViewSize()
	this.ViewWidth--

	this.Prompt = prompt_
	if this.Prompt != nil {
		this.ViewWidth = this.ViewWidth - this.Prompt()
	}
	for {
		stdOut.Flush()
		shineCursor()
		this.Unicode, this.Keycode, this.ShiftState = conio.GetKey()
		var f KeyFuncT
		var ok bool
		if (this.ShiftState & conio.ALT_PRESSED) != 0 {
			f, ok = altMap[this.Keycode]
			if !ok {
				continue
			}
		} else if this.Unicode != 0 {
			f, ok = keyMap[this.Unicode]
			if !ok {
				//f = KeyFuncInsertReport
				f = &KeyGoFuncT{KeyFuncInsertSelf}
			}
		} else {
			f, ok = scanMap[this.Keycode]
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
