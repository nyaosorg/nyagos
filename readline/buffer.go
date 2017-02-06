package readline

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"unicode"

	. "../conio"
	"github.com/mattn/go-colorable"
)

var Console = colorable.NewColorableStdout()

var hasCache = map[rune]bool{}

func PutRune(ch rune) {
	if hasCache[ch] {
		fmt.Fprintf(Console, "%c", ch)
	} else {
		pre_x, pre_y := GetLocate()
		fmt.Fprintf(Console, "%c", ch)
		post_x, post_y := GetLocate()
		if post_y == pre_y && post_x > pre_x {
			hasCache[ch] = true
			SetCharWidth(ch, post_x-pre_x)
		}
	}
}

func PutRunes(ch rune, n int) {
	if n <= 0 {
		return
	}
	PutRune(ch)
	for i := 1; i < n; i++ {
		fmt.Fprintf(Console, "%c", ch)
	}
}

func Backspace(n int) {
	if n > 1 {
		fmt.Fprintf(Console, "\x1B[%dC", n)
	} else if n == 1 {
		fmt.Fprint(Console, "\b")
	}
}

func Eraseline() {
	fmt.Fprint(Console, "\x1B[0K")
}

const FORBIDDEN_WIDTH = 3 // = lastcolumn(1) and FULLWIDTHCHAR-SIZE(2)

type Buffer struct {
	Buffer         []rune
	Length         int
	Cursor         int
	Unicode        rune
	Keycode        uint16
	ShiftState     uint32
	ViewStart      int
	TermWidth      int // == TopColumn + ViewWidth + FORBIDDEN_WIDTH
	TopColumn      int // == width of Prompt
	HistoryPointer int
	Context        context.Context
	History        IHistory
	Prompt         func() (int, error)
}

func (this *Buffer) ViewWidth() int {
	return this.TermWidth - this.TopColumn - FORBIDDEN_WIDTH
}

func (this *Buffer) Insert(pos int, c []rune) {
	n := len(c)
	for this.Length+n >= len(this.Buffer) {
		tmp := make([]rune, len(this.Buffer)*2)
		copy(tmp, this.Buffer)
		this.Buffer = tmp
	}
	for i := this.Length - 1; i >= pos; i-- {
		this.Buffer[i+n] = this.Buffer[i]
	}
	for i := 0; i < n; i++ {
		this.Buffer[pos+i] = c[i]
	}
	this.Length += n
}

// Insert String :s at :pos
// returns
//    count of rune
func (this *Buffer) InsertString(pos int, s string) int {
	list := []rune(s)
	this.Insert(pos, list)
	return len(list)
}

func (this *Buffer) Delete(pos int, n int) int {
	if n <= 0 || this.Length < pos+n {
		return 0
	}
	delw := this.GetWidthBetween(pos, pos+n)
	for i := pos; i < this.Length-n; i++ {
		this.Buffer[i] = this.Buffer[i+n]
	}
	this.Length -= n
	return delw
}

func (this *Buffer) InsertAndRepaint(str string) {
	this.ReplaceAndRepaint(this.Cursor, str)
}

func (this *Buffer) ResetViewStart() {
	this.ViewStart = 0
	w := 0
	for i := 0; i <= this.Cursor; i++ {
		w += GetCharWidth(this.Buffer[i])
		for w >= this.ViewWidth() {
			w -= GetCharWidth(this.Buffer[this.ViewStart])
			this.ViewStart++
		}
	}
}

func (this *Buffer) ReplaceAndRepaint(pos int, str string) {
	// Cursor rewind
	Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))

	// Replace Buffer
	this.Delete(pos, this.Cursor-pos)

	// Define ViewStart , Cursor
	this.Cursor = pos + this.InsertString(pos, str)
	this.ResetViewStart()

	// Repaint
	w := 0
	for i := this.ViewStart; i < this.Cursor; i++ {
		PutRune(this.Buffer[i])
		w += GetCharWidth(this.Buffer[i])
	}
	bs := 0
	for i := this.Cursor; i < this.Length; i++ {
		w1 := GetCharWidth(this.Buffer[i])
		if w+w1 >= this.ViewWidth() {
			break
		}
		PutRune(this.Buffer[i])
		w += w1
		bs += w1
	}
	Eraseline()
	if bs > 0 {
		Backspace(bs)
	}
}

func (this *Buffer) GetWidthBetween(from int, to int) int {
	width := 0
	for i := from; i < to; i++ {
		width += GetCharWidth(this.Buffer[i])
	}
	return width
}

func (this *Buffer) Repaint(pos int, del int) {
	bs := 0
	vp := this.GetWidthBetween(this.ViewStart, pos)

	for i := pos; i < this.Length; i++ {
		w1 := GetCharWidth(this.Buffer[i])
		if vp+w1 >= this.ViewWidth() {
			break
		}
		PutRune(this.Buffer[i])
		vp += w1
		bs += w1
	}
	Eraseline()
	if del > 0 {
		Backspace(bs)
	} else {
		// for readline_keyfunc.go: KeyFuncInsertSelf()
		Backspace(bs + del)
	}
}

func (this *Buffer) RepaintAfterPrompt() {
	this.ResetViewStart()
	for i := this.ViewStart; i < this.Cursor; i++ {
		PutRune(this.Buffer[i])
	}
	this.Repaint(this.Cursor, 0)
}

func (this *Buffer) RepaintAll() {
	this.TopColumn, _ = this.Prompt()
	this.RepaintAfterPrompt()
}

func (this Buffer) String() string {
	var result bytes.Buffer
	for i := 0; i < this.Length; i++ {
		result.WriteRune(this.Buffer[i])
	}
	return result.String()
}

var Delimiters = "\"'"

func (this *Buffer) CurrentWordTop() (wordTop int) {
	wordTop = -1
	quotedchar := '\000'
	for i := 0; i < this.Cursor; i++ {
		if quotedchar == '\000' {
			if strings.ContainsRune(Delimiters, this.Buffer[i]) {
				quotedchar = this.Buffer[i]
			}
		} else if this.Buffer[i] == quotedchar {
			quotedchar = '\000'
		}
		if unicode.IsSpace(this.Buffer[i]) && quotedchar == '\000' {
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

func (this *Buffer) CurrentWord() (string, int) {
	var buffer bytes.Buffer
	start := this.CurrentWordTop()
	for i := start; i < this.Cursor; i++ {
		buffer.WriteRune(this.Buffer[i])
	}
	return buffer.String(), start
}
