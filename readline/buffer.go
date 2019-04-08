package readline

import (
	"fmt"
	"io"
	"strings"
	"unicode"

	"github.com/mattn/go-runewidth"
	"github.com/mattn/go-tty"
)

var SurrogatePairOk = false

func (this *Buffer) PutRune(ch rune) {
	if ch < ' ' {
		this.Out.WriteByte('^')
		this.Out.WriteByte(byte('A' + (ch - 1)))
	} else if (ch >= 0x10000 && !SurrogatePairOk) || runewidth.RuneWidth(ch) == 0 {
		fmt.Fprintf(this.Out, "<%X>", ch)
	} else {
		this.Out.WriteRune(ch)
	}
}

func (this *Buffer) PutRunes(ch rune, n int) {
	if n <= 0 {
		return
	}
	this.PutRune(ch)
	for i := 1; i < n; i++ {
		this.Out.WriteRune(ch)
	}
}

func (this *Buffer) Backspace(n int) {
	if n > 1 {
		fmt.Fprintf(this.Out, "\x1B[%dD", n)
	} else if n == 1 {
		this.Out.WriteByte('\b')
	}
}

func (this *Buffer) Eraseline() {
	io.WriteString(this.Out, "\x1B[0K")
}

const forbiddenWidth = 3 // = lastcolumn(1) and FULLWIDTHCHAR-SIZE(2)

type Buffer struct {
	*Editor
	Buffer         []rune
	TTY            *tty.TTY
	ViewStart      int
	TermWidth      int // == TopColumn + ViewWidth + forbiddenWidth
	TopColumn      int // == width of Prompt
	HistoryPointer int
}

func (this *Buffer) ViewWidth() int {
	return this.TermWidth - this.TopColumn - forbiddenWidth
}

func (this *Buffer) insert(csrPos int, insStr []rune) {
	// expand buffer
	this.Buffer = append(this.Buffer, insStr...)

	// shift original string to make area
	copy(this.Buffer[csrPos+len(insStr):], this.Buffer[csrPos:])

	// insert insStr
	copy(this.Buffer[csrPos:csrPos+len(insStr)], insStr)
}

// Insert String :s at :pos (Do not update screen)
// returns
//    count of rune
func (this *Buffer) InsertString(pos int, s string) int {
	list := []rune(s)
	this.insert(pos, list)
	return len(list)
}

func (this *Buffer) Delete(pos int, n int) int {
	if n <= 0 || len(this.Buffer) < pos+n {
		return 0
	}
	delw := this.GetWidthBetween(pos, pos+n)
	copy(this.Buffer[pos:], this.Buffer[pos+n:])
	this.Buffer = this.Buffer[:len(this.Buffer)-n]
	return delw
}

func (this *Buffer) InsertAndRepaint(str string) {
	this.ReplaceAndRepaint(this.Cursor, str)
}

func (this *Buffer) ResetViewStart() {
	this.ViewStart = 0
	w := 0
	for i := 0; i <= this.Cursor && i < len(this.Buffer); i++ {
		w += GetCharWidth(this.Buffer[i])
		for w >= this.ViewWidth() {
			if this.ViewStart >= len(this.Buffer) {
				// When standard output is redirected.
				return
			}
			w -= GetCharWidth(this.Buffer[this.ViewStart])
			this.ViewStart++
		}
	}
}

func (this *Buffer) ReplaceAndRepaint(pos int, str string) {
	// Cursor rewind
	this.Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))

	// Replace Buffer
	this.Delete(pos, this.Cursor-pos)

	// Define ViewStart , Cursor
	this.Cursor = pos + this.InsertString(pos, str)
	this.ResetViewStart()

	// Repaint
	w := 0
	for i := this.ViewStart; i < this.Cursor; i++ {
		this.PutRune(this.Buffer[i])
		w += GetCharWidth(this.Buffer[i])
	}
	bs := 0
	for i := this.Cursor; i < len(this.Buffer); i++ {
		w1 := GetCharWidth(this.Buffer[i])
		if w+w1 >= this.ViewWidth() {
			break
		}
		this.PutRune(this.Buffer[i])
		w += w1
		bs += w1
	}
	this.Eraseline()
	if bs > 0 {
		this.Backspace(bs)
	}
}

func (this *Buffer) GetWidthBetween(from int, to int) int {
	width := 0
	for i := from; i < to; i++ {
		width += GetCharWidth(this.Buffer[i])
	}
	return width
}

// Repaint buffer[pos:] + " \b"*del but do not rewind cursor position
func (this *Buffer) Repaint(pos int, del int) {
	bs := 0
	vp := this.GetWidthBetween(this.ViewStart, pos)

	for i := pos; i < len(this.Buffer); i++ {
		w1 := GetCharWidth(this.Buffer[i])
		if vp+w1 >= this.ViewWidth() {
			break
		}
		this.PutRune(this.Buffer[i])
		vp += w1
		bs += w1
	}
	this.Eraseline()
	if del > 0 {
		this.Backspace(bs)
	} else {
		// for readline_keyfunc.go: KeyFuncInsertSelf()
		this.Backspace(bs + del)
	}
}

func (this *Buffer) RepaintAfterPrompt() {
	this.ResetViewStart()
	for i := this.ViewStart; i < this.Cursor; i++ {
		this.PutRune(this.Buffer[i])
	}
	this.Repaint(this.Cursor, 0)
}

func (this *Buffer) RepaintAll() {
	this.Out.Flush()
	this.TopColumn, _ = this.Prompt()
	this.RepaintAfterPrompt()
}

func (this *Buffer) SubString(start, end int) string {
	var result strings.Builder
	for i := start; i < end; i++ {
		result.WriteRune(this.Buffer[i])
	}
	return result.String()
}

func (this Buffer) String() string {
	return this.SubString(0, len(this.Buffer))
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
	start := this.CurrentWordTop()
	return this.SubString(start, this.Cursor), start
}
