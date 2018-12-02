package readline

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/mattn/go-tty"
)

var FlushBeforeReadline = false

type Result int

const (
	CONTINUE Result = iota
	ENTER    Result = iota
	ABORT    Result = iota
	INTR     Result = iota
)

func (this Result) String() string {
	switch this {
	case CONTINUE:
		return "CONTINUE"
	case ENTER:
		return "ENTER"
	case ABORT:
		return "ABORT"
	case INTR:
		return "INTR"
	default:
		return "ERROR"
	}
}

type KeyFuncT interface {
	Call(ctx context.Context, buffer *Buffer) Result
}

type KeyGoFuncT struct {
	Func func(ctx context.Context, buffer *Buffer) Result
	Name string
}

func (this *KeyGoFuncT) Call(ctx context.Context, buffer *Buffer) Result {
	if this.Func == nil {
		return CONTINUE
	}
	return this.Func(ctx, buffer)
}

func (this KeyGoFuncT) String() string {
	return this.Name
}

var keyMap = map[string]KeyFuncT{
	name2char[K_CTRL_A]:    name2func(F_BEGINNING_OF_LINE),
	name2char[K_CTRL_B]:    name2func(F_BACKWARD_CHAR),
	name2char[K_BACKSPACE]: name2func(F_BACKWARD_DELETE_CHAR),
	name2char[K_CTRL_C]:    name2func(F_INTR),
	name2char[K_CTRL_D]:    name2func(F_DELETE_OR_ABORT),
	name2char[K_CTRL_E]:    name2func(F_END_OF_LINE),
	name2char[K_CTRL_F]:    name2func(F_FORWARD_CHAR),
	name2char[K_CTRL_H]:    name2func(F_BACKWARD_DELETE_CHAR),
	name2char[K_CTRL_K]:    name2func(F_KILL_LINE),
	name2char[K_CTRL_L]:    name2func(F_CLEAR_SCREEN),
	name2char[K_CTRL_M]:    name2func(F_ACCEPT_LINE),
	name2char[K_CTRL_R]:    name2func(F_ISEARCH_BACKWARD),
	name2char[K_CTRL_U]:    name2func(F_UNIX_LINE_DISCARD),
	name2char[K_CTRL_Y]:    name2func(F_YANK),
	name2char[K_DELETE]:    name2func(F_DELETE_CHAR),
	name2char[K_ENTER]:     name2func(F_ACCEPT_LINE),
	name2char[K_ESCAPE]:    name2func(F_KILL_WHOLE_LINE),
	name2char[K_CTRL_N]:    name2func(F_HISTORY_DOWN),
	name2char[K_CTRL_P]:    name2func(F_HISTORY_UP),
	name2char[K_CTRL_Q]:    name2func(F_QUOTED_INSERT),
	name2char[K_CTRL_T]:    name2func(F_SWAPCHAR),
	name2char[K_CTRL_V]:    name2func(F_QUOTED_INSERT),
	name2char[K_CTRL_W]:    name2func(F_UNIX_WORD_RUBOUT),
	name2char[K_CTRL]:      name2func(F_PASS),
	name2char[K_DELETE]:    name2func(F_DELETE_CHAR),
	name2char[K_END]:       name2func(F_END_OF_LINE),
	name2char[K_HOME]:      name2func(F_BEGINNING_OF_LINE),
	name2char[K_LEFT]:      name2func(F_BACKWARD_CHAR),
	name2char[K_RIGHT]:     name2func(F_FORWARD_CHAR),
	name2char[K_SHIFT]:     name2func(F_PASS),
	name2char[K_DOWN]:      name2func(F_HISTORY_DOWN),
	name2char[K_UP]:        name2func(F_HISTORY_UP),
	name2char[K_ALT_V]:     name2func(F_YANK),
	name2char[K_ALT_Y]:     name2func(F_YANK_WITH_QUOTE),
}

func normWord(src string) string {
	return strings.Replace(strings.ToUpper(src), "-", "_", -1)
}

func BindKeyFunc(keyName string, funcValue KeyFuncT) error {
	keyName_ := normWord(keyName)
	if charValue, charOk := name2char[keyName_]; charOk {
		keyMap[charValue] = funcValue
		return nil
	}
	return fmt.Errorf("%s: no such keyname", keyName)
}

func BindKeyClosure(name string, f func(context.Context, *Buffer) Result) error {
	return BindKeyFunc(name, &KeyGoFuncT{Func: f, Name: "annonymous"})
}

func GetBindKey(keyName string) KeyFuncT {
	keyName_ := normWord(keyName)
	if charValue, charOk := name2char[keyName_]; charOk {
		return keyMap[charValue]
	} else {
		return nil
	}
}

func GetFunc(funcName string) (KeyFuncT, error) {
	rc := name2func(normWord(funcName))
	if rc != nil {
		return rc, nil
	} else {
		return nil, fmt.Errorf("%s: not found in the function-list", funcName)
	}
}

func BindKeySymbol(keyName, funcName string) error {
	funcValue := name2func(normWord(funcName))
	if funcValue == nil {
		return fmt.Errorf("%s: no such function.", funcName)
	}
	return BindKeyFunc(keyName, funcValue)
}

type EmptyHistory struct{}

func (this *EmptyHistory) Len() int      { return 0 }
func (this *EmptyHistory) At(int) string { return "" }

const (
	CURSOR_OFF = "\x1B[?25l"
	CURSOR_ON  = "\x1B[?25h\x1B[s\x1B[u"
)

var CtrlC = errors.New("^C")

func getKey(tty1 *tty.TTY) (string, error) {
	var buffer strings.Builder
	for {
		r, err := tty1.ReadRune()
		if err != nil {
			return "", err
		}
		if r == 0 {
			continue
		}
		buffer.WriteRune(r)
		if !tty1.Buffered() {
			return buffer.String(), nil
		}
	}
}

var mu sync.Mutex

// Call LineEditor
// - ENTER typed -> returns TEXT and nil
// - CTRL-C typed -> returns "" and readline.CtrlC
// - CTRL-D typed -> returns "" and io.EOF
func (session *Editor) ReadLine(ctx context.Context) (string, error) {
	if session.Writer == nil {
		panic("readline.Editor.Writer is not set. Set an instance such as go-colorable.NewColorableStdout()")
	}
	defer func() {
		session.Writer.WriteString(CURSOR_ON)
		session.Writer.Flush()
	}()

	if session.Prompt == nil {
		session.Prompt = func() (int, error) {
			session.Writer.WriteString("\n> ")
			return 2, nil
		}
	}
	if session.History == nil {
		session.History = new(EmptyHistory)
	}
	this := Buffer{
		Editor:         session,
		Buffer:         make([]rune, 20),
		HistoryPointer: session.History.Len(),
	}

	tty1, err := tty.Open()
	if err != nil {
		return "", fmt.Errorf("go-tty.Open: %s", err.Error())
	}
	this.TTY = tty1
	defer tty1.Close()

	this.TermWidth, _, err = tty1.Size()
	if err != nil {
		return "", fmt.Errorf("go-tty.Size: %s", err.Error())
	}

	var err1 error
	this.TopColumn, err1 = session.Prompt()
	if err1 != nil {
		// unable to get prompt-string.
		fmt.Fprintf(this.Writer, "%s\n$ ", err1.Error())
		this.TopColumn = 2
	} else if this.TopColumn >= this.TermWidth-3 {
		// ViewWidth is too narrow to edit.
		io.WriteString(this.Writer, "\n")
		this.TopColumn = 0
	}
	this.InsertString(0, session.Default)
	if this.Cursor > this.Length {
		this.Cursor = this.Length
	}
	this.RepaintAfterPrompt()

	cursorOnSwitch := false

	ws := tty1.SIGWINCH()
	go func(lastw int) {
		for ws1 := range ws {
			w := ws1.W
			if lastw != w {
				mu.Lock()
				this.TermWidth = w
				fmt.Fprintf(this.Writer, "\x1B[%dG", this.TopColumn+1)
				this.RepaintAfterPrompt()
				mu.Unlock()
				lastw = w
			}
		}
	}(this.TermWidth)

	for {
		mu.Lock()
		if !cursorOnSwitch {
			io.WriteString(this.Writer, CURSOR_ON)
			cursorOnSwitch = true
		}
		this.Writer.Flush()

		mu.Unlock()
		key1, err := getKey(tty1)
		if err != nil {
			return "", err
		}
		mu.Lock()
		f, ok := keyMap[key1]
		if !ok {
			f = &KeyGoFuncT{
				Func: func(ctx context.Context, this *Buffer) Result {
					return KeyFuncInsertSelf(ctx, this, key1)
				},
				Name: key1,
			}
		}

		if fg, ok := f.(*KeyGoFuncT); !ok || fg.Func != nil {
			io.WriteString(this.Writer, CURSOR_OFF)
			cursorOnSwitch = false
		}
		rc := f.Call(ctx, &this)
		if rc != CONTINUE {
			this.Writer.WriteByte('\n')
			if !cursorOnSwitch {
				io.WriteString(this.Writer, CURSOR_ON)
			}
			this.Writer.Flush()
			result := this.String()
			mu.Unlock()
			if rc == ENTER {
				return result, nil
			} else if rc == INTR {
				return result, CtrlC
			} else {
				return result, io.EOF
			}
		}
		mu.Unlock()
	}
}
