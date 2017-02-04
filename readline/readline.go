package readline

import (
	"fmt"
	"strings"

	. "../conio"
	"github.com/zetamatta/go-getch"
)

type Result int

const (
	CONTINUE Result = iota
	ENTER    Result = iota
	ABORT    Result = iota
)

func (this Result) String() string {
	switch this {
	case CONTINUE:
		return "CONTINUE"
	case ENTER:
		return "ENTER"
	case ABORT:
		return "ABORT"
	default:
		return "ERROR"
	}
}

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
	name2char[K_CTRL_R]: name2func[F_ISEARCH_BACKWARD],
	name2char[K_CTRL_U]: name2func[F_UNIX_LINE_DISCARD],
	name2char[K_CTRL_Y]: name2func[F_YANK],
	name2char[K_DELETE]: name2func[F_DELETE_CHAR],
	name2char[K_ENTER]:  name2func[F_ACCEPT_LINE],
	name2char[K_ESCAPE]: name2func[F_KILL_WHOLE_LINE],
	name2char[K_CTRL_N]: name2func[F_HISTORY_DOWN],
	name2char[K_CTRL_P]: name2func[F_HISTORY_UP],
	name2char[K_CTRL_T]: name2func[F_SWAPCHAR],
	name2char[K_CTRL_W]: name2func[F_UNIX_WORD_RUBOUT],
}

var scanMap = map[uint16]KeyFuncT{
	name2scan[K_CTRL]:   name2func[F_PASS],
	name2scan[K_DELETE]: name2func[F_DELETE_CHAR],
	name2scan[K_END]:    name2func[F_END_OF_LINE],
	name2scan[K_HOME]:   name2func[F_BEGINNING_OF_LINE],
	name2scan[K_LEFT]:   name2func[F_BACKWARD_CHAR],
	name2scan[K_RIGHT]:  name2func[F_FORARD_CHAR],
	name2scan[K_SHIFT]:  name2func[F_PASS],
	name2scan[K_DOWN]:   name2func[F_HISTORY_DOWN],
	name2scan[K_UP]:     name2func[F_HISTORY_UP],
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

type LineEditorAbort struct{}

func (this *LineEditorAbort) Error() string {
	return "^D"
}

// Call LineEditor
// - ENTER typed -> returns TEXT and nil
// - CTRL-C typed -> returns "" and nil
// - CTRL-D typed -> returns "" and LineEditorAbort
func (session *LineEditor) ReadLine() (string, error) {
	this := Buffer{
		Buffer:  make([]rune, 20),
		Session: session,
	}
	this.TermWidth, _ = GetScreenBufferInfo().ViewSize()

	var err1 error
	this.TopColumn, err1 = session.Prompt(session)
	if err1 != nil {
		// unable to get prompt-string.
		fmt.Fprintf(stdOut, "%s\n$ ", err1.Error())
		this.TopColumn = 2
	} else if this.TopColumn >= this.TermWidth-3 {
		// ViewWidth is too narrow to edit.
		fmt.Fprint(stdOut, "\n")
		this.TopColumn = 0
	}
	for {
		stdOut.Flush()
		shineCursor()
		e := getch.All()
		if e.Resize != nil {
			w := int(e.Resize.Width)
			if this.TermWidth != w {
				this.TermWidth = w
				fmt.Fprintf(stdOut, "\x1B[%dG", this.TopColumn+1)
				stdOut.Flush()
				shineCursor()
				this.RepaintAfterPrompt()
				stdOut.Flush()
				shineCursor()
			}
			continue
		}
		if e.Key == nil {
			continue
		}
		this.Unicode = e.Key.Rune
		this.Keycode = e.Key.Scan
		this.ShiftState = e.Key.Shift
		var f KeyFuncT
		var ok bool
		if (this.ShiftState&getch.ALT_PRESSED) != 0 &&
			(this.ShiftState&getch.CTRL_PRESSED) == 0 {
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
			result := this.String()
			if result == "" {
				session.Pointer = session.History.Len()
			}
			if last := session.History.At(-1); result != last {
				session.History.Push(result)
			}
			if rc == ENTER {
				return result, nil
			} else {
				return result, new(LineEditorAbort)
			}
		}
	}
}
