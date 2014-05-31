package conio

import "bufio"
import "bytes"
import "fmt"
import "os"
import "unicode"

import "github.com/mattn/go-runewidth"

var widthCache = make(map[rune]int)

func getCharWidth(n rune) int {
	width, ok := widthCache[n]
	if !ok {
		width = runewidth.RuneWidth(n)
		widthCache[n] = width
	}
	return width
	// if n > 0xFF {
	//	return 2;
	//}else{
	//	return 1;
	//}
}

var stdOut *bufio.Writer = bufio.NewWriter(os.Stdout)

func PutRep(ch rune, n int) {
	for i := 0; i < n; i++ {
		stdOut.WriteRune(ch)
	}
}

func Backspace(n int) {
	stdOut.Flush()
	x, y := GetLocate()
	Locate(x-n, y)
}

func shineCursor() {
	x, y := GetLocate()
	Locate(x, y)
}

type ReadLineBuffer struct {
	Buffer    []rune
	Length    int
	Cursor    int
	Unicode   rune
	Keycode   uint16
	ViewStart int
	ViewWidth int
}

func (this *ReadLineBuffer) Insert(pos int, c []rune) bool {
	n := len(c)
	for this.Length+n >= len(this.Buffer) {
		tmp := make([]rune, len(this.Buffer)*2)
		copy(tmp, this.Buffer)
		this.Buffer = tmp
	}
	for i := this.Length; i >= pos; i-- {
		this.Buffer[i+n] = this.Buffer[i]
	}
	for i := 0; i < n; i++ {
		this.Buffer[pos+i] = c[i]
	}
	this.Length += n
	return true
}

func (this *ReadLineBuffer) InsertString(pos int, s string) int {
	list := make([]rune, 0)
	for _, r := range s {
		list = append(list, r)
	}
	if this.Insert(pos, list) {
		return len(list)
	} else {
		return -1
	}
}

func (this *ReadLineBuffer) Delete(pos int, n int) int {
	if this.Length < pos+n {
		return 0
	}
	delw := 0
	for i := pos; i < pos+n; i++ {
		delw += getCharWidth(this.Buffer[i])
	}
	for i := pos; i < this.Length-n; i++ {
		this.Buffer[i] = this.Buffer[i+n]
	}
	this.Length -= n
	return delw
}

func (this *ReadLineBuffer) ReplaceAndRepaint(pos int, str string) {
	n := this.Cursor - pos
	this.Delete(pos, n)
	this.InsertString(pos, str)
	if pos < this.ViewStart {
		Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
	} else {
		Backspace(this.GetWidthBetween(pos, this.Cursor))
	}
	this.Cursor = pos
	for _, ch := range str {
		if this.Cursor >= this.ViewStart {
			PutRep(ch, 1)
		}
		this.Cursor++
	}
	this.Repaint(this.Cursor, 0)
}

func (this *ReadLineBuffer) GetWidthBetween(from int, to int) int {
	width := 0
	for i := from; i < to; i++ {
		width += getCharWidth(this.Buffer[i])
	}
	return width
}

func (this *ReadLineBuffer) Repaint(pos int, del int) {
	bs := 0
	vp := this.GetWidthBetween(this.ViewStart, pos)

	for i := pos; i < this.Length; i++ {
		w1 := getCharWidth(this.Buffer[i])
		vp += w1
		if vp >= this.ViewWidth {
			break
		}
		PutRep(this.Buffer[i], 1)
		bs += w1
	}
	PutRep(' ', del)
	Backspace(bs + del)
}

func (this *ReadLineBuffer) RepaintAll(header string) {
	PutRep('\r', 1)
	for _, ch := range header {
		PutRep(ch, 1)
	}
	for i := this.ViewStart; i < this.Cursor; i++ {
		PutRep(this.Buffer[i], 1)
	}
	this.Repaint(this.Cursor, 0)
}

func (this ReadLineBuffer) String() string {
	var result bytes.Buffer
	for i := 0; i < this.Length; i++ {
		result.WriteRune(this.Buffer[i])
	}
	return result.String()
}

func (this *ReadLineBuffer) CurrentWordTop() (wordTop int) {
	wordTop = -1
	isQuoted := false
	for i := 0; i < this.Cursor; i++ {
		if this.Buffer[i] == '"' {
			isQuoted = !isQuoted
		}
		if unicode.IsSpace(this.Buffer[i]) && !isQuoted {
			wordTop = -1
		} else if wordTop < 0 {
			wordTop = i
		}
	}
	if wordTop < 0 {
		return this.Cursor
	} else {
		return wordTop
	}
}

func (this *ReadLineBuffer) CurrentWord() (string, int) {
	var buffer bytes.Buffer
	start := this.CurrentWordTop()
	for i := start; i < this.Cursor; i++ {
		if this.Buffer[i] != '"' {
			buffer.WriteRune(this.Buffer[i])
		}
	}
	return buffer.String(), start
}

type KeyFuncResult int

const (
	CONTINUE KeyFuncResult = iota
	ENTER    KeyFuncResult = iota
	ABORT    KeyFuncResult = iota
)

func KeyFuncPass(this *ReadLineBuffer) KeyFuncResult {
	return CONTINUE
}

func KeyFuncEnter(this *ReadLineBuffer) KeyFuncResult { // Ctrl-M
	return ENTER
}

func KeyFuncHead(this *ReadLineBuffer) KeyFuncResult { // Ctrl-A
	Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
	this.Cursor = 0
	this.ViewStart = 0
	this.Repaint(0, 1)
	return CONTINUE
}

func KeyFuncBackword(this *ReadLineBuffer) KeyFuncResult { // Ctrl-B
	if this.Cursor <= 0 {
		return CONTINUE
	}
	this.Cursor--
	if this.Cursor < this.ViewStart {
		this.ViewStart--
		this.Repaint(this.Cursor, 1)
	} else {
		Backspace(getCharWidth(this.Buffer[this.Cursor]))
	}
	return CONTINUE
}

func KeyFuncTail(this *ReadLineBuffer) KeyFuncResult { // Ctrl-E
	allength := this.GetWidthBetween(this.ViewStart, this.Length)
	if allength < this.ViewWidth {
		for ; this.Cursor < this.Length; this.Cursor++ {
			PutRep(this.Buffer[this.Cursor], 1)
		}
	} else {
		PutRep('\a', 1)
		Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
		this.ViewStart = this.Length - 1
		w := getCharWidth(this.Buffer[this.ViewStart])
		for {
			if this.ViewStart <= 0 {
				break
			}
			w_ := w + getCharWidth(this.Buffer[this.ViewStart-1])
			if w_ >= this.ViewWidth {
				break
			}
			w = w_
			this.ViewStart--
		}
		for this.Cursor = this.ViewStart; this.Cursor < this.Length; this.Cursor++ {
			PutRep(this.Buffer[this.Cursor], 1)
		}
	}
	return CONTINUE
}

func KeyFuncForward(this *ReadLineBuffer) KeyFuncResult { // Ctrl-F
	if this.Cursor >= this.Length {
		return CONTINUE
	}
	w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
	if w < this.ViewWidth {
		// No Scroll
		PutRep(this.Buffer[this.Cursor], 1)
	} else {
		// Right Scroll
		Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
		if getCharWidth(this.Buffer[this.Cursor]) > getCharWidth(this.Buffer[this.ViewStart]) {
			this.ViewStart++
		}
		this.ViewStart++
		for i := this.ViewStart; i <= this.Cursor; i++ {
			PutRep(this.Buffer[i], 1)
		}
		PutRep(' ', 1)
		Backspace(1)
	}
	this.Cursor++
	return CONTINUE
}

func KeyFuncBackSpace(this *ReadLineBuffer) KeyFuncResult { // Backspace
	if this.Cursor > 0 {
		this.Cursor--
		delw := this.Delete(this.Cursor, 1)
		if this.Cursor >= this.ViewStart {
			Backspace(delw)
		} else {
			this.ViewStart = this.Cursor
		}
		this.Repaint(this.Cursor, delw)
	}
	return CONTINUE
}

func KeyFuncDelete(this *ReadLineBuffer) KeyFuncResult { // Del
	delw := this.Delete(this.Cursor, 1)
	this.Repaint(this.Cursor, delw)
	return CONTINUE
}

func KeyFuncDeleteOrAbort(this *ReadLineBuffer) KeyFuncResult { // Ctrl-D
	if this.Length > 0 {
		return KeyFuncDelete(this)
	} else {
		return ABORT
	}
}

func KeyFuncInsertSelf(this *ReadLineBuffer) KeyFuncResult {
	ch := this.Unicode
	if ch < 0x20 || !this.Insert(this.Cursor, []rune{ch}) {
		return CONTINUE
	}
	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	w1 := getCharWidth(ch)
	if w+w1 >= this.ViewWidth {
		// scroll left
		Backspace(w)
		if getCharWidth(this.Buffer[this.ViewStart]) < w1 {
			this.ViewStart++
		}
		this.ViewStart++
		for i := this.ViewStart; i <= this.Cursor; i++ {
			PutRep(this.Buffer[i], 1)
		}
		PutRep(' ', 1)
		Backspace(1)
	} else {
		this.Repaint(this.Cursor, -w1)
	}
	this.Cursor++
	return CONTINUE
}

func KeyFuncInsertReport(this *ReadLineBuffer) KeyFuncResult {
	L := this.InsertString(this.Cursor, fmt.Sprintf("[%X]", this.Unicode))
	if L >= 0 {
		this.Repaint(this.Cursor, -L)
		this.Cursor += L
	}
	return CONTINUE
}

func KeyFuncClearAfter(this *ReadLineBuffer) KeyFuncResult {
	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	i := this.Cursor
	bs := 0
	for i < this.Length && w < this.ViewWidth {
		w1 := getCharWidth(this.Buffer[i])
		PutRep(' ', w1)
		i++
		w += w1
		bs += w1
	}
	Backspace(bs)
	this.Length = this.Cursor
	return CONTINUE
}

func KeyFuncClear(this *ReadLineBuffer) KeyFuncResult {
	KeyFuncClearAfter(this)
	width := this.GetWidthBetween(this.ViewStart, this.Cursor)
	Backspace(width)
	PutRep(' ', width)
	Backspace(width)
	this.Length = 0
	this.Cursor = 0
	this.ViewStart = 0
	return CONTINUE
}

func KeyFuncClearBefore(this *ReadLineBuffer) KeyFuncResult {
	keta := this.Delete(0, this.Cursor)
	Backspace(keta)
	this.Cursor = 0
	this.Repaint(0, keta)
	return CONTINUE
}

var KeyMap = map[rune]func(*ReadLineBuffer) KeyFuncResult{
	'\r':         KeyFuncEnter,
	'\x01':       KeyFuncHead,
	'\x02':       KeyFuncBackword,
	'\x05':       KeyFuncTail,
	'\x06':       KeyFuncForward,
	'\b':         KeyFuncBackSpace,
	'\x04':       KeyFuncDeleteOrAbort,
	'\x7F':       KeyFuncDelete,
	('K' & 0x1F): KeyFuncClearAfter,
	'\x1B':       KeyFuncClear,
	('U' & 0x1F): KeyFuncClearBefore,
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

func ReadLine() (string, KeyFuncResult) {
	var this ReadLineBuffer
	this.Buffer = make([]rune, 20)
	this.Length = 0
	this.Cursor = 0
	this.ViewStart = 0
	this.ViewWidth = 60
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
