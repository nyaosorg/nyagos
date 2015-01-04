package conio

import (
	"bufio"
	"bytes"
	"os"
	"unicode"
)

var stdOut *bufio.Writer = bufio.NewWriter(os.Stdout)

var hasCache = map[rune]bool{}

func PutRep(ch rune, n int) {
	if n <= 0 {
		return
	}
	if hasCache[ch] {
		stdOut.WriteRune(ch)
	} else {
		stdOut.Flush()
		pre_x, pre_y := GetLocate()
		stdOut.WriteRune(ch)
		stdOut.Flush()
		post_x, post_y := GetLocate()
		if post_y == pre_y && post_x > pre_x {
			hasCache[ch] = true
			SetCharWidth(ch, post_x-pre_x)
		}
	}
	for i := 1; i < n; i++ {
		stdOut.WriteRune(ch)
	}
}

func Backspace(n int) {
	if n > 0 {
		stdOut.Flush()
		x, y := GetLocate()
		Locate(x-n, y)
	}
}

func shineCursor() {
	x, y := GetLocate()
	Locate(x, y)
}

type Buffer struct {
	Buffer     []rune
	Length     int
	Cursor     int
	Unicode    rune
	Keycode    uint16
	ShiftState uint32
	ViewStart  int
	ViewWidth  int
	Prompt     func() int
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

func (this *Buffer) ReplaceAndRepaint(pos int, str string) {
	// Cursor rewind
	Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))

	// Replace Buffer
	this.Delete(pos, this.Cursor-pos)

	// Define ViewStart , Cursor
	this.Cursor = pos + this.InsertString(pos, str)
	this.ViewStart = 0
	w := 0
	for i := 0; i < this.Cursor; i++ {
		w += GetCharWidth(this.Buffer[i])
		for w >= this.ViewWidth {
			w -= GetCharWidth(this.Buffer[this.ViewStart])
			this.ViewStart++
		}
	}

	// Repaint
	w = 0
	for i := this.ViewStart; i < this.Cursor; i++ {
		PutRep(this.Buffer[i], 1)
		w += GetCharWidth(this.Buffer[i])
	}
	bs := 0
	for i := this.Cursor; i < this.Length; i++ {
		w1 := GetCharWidth(this.Buffer[i])
		if w+w1 >= this.ViewWidth {
			break
		}
		PutRep(this.Buffer[i], 1)
		w += w1
		bs += w1
	}
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

func (this *Buffer) RepaintAll() {
	if this.Prompt != nil {
		this.Prompt()
	}
	for i := this.ViewStart; i < this.Cursor; i++ {
		PutRep(this.Buffer[i], 1)
	}
	this.Repaint(this.Cursor, 0)
}

func (this Buffer) String() string {
	var result bytes.Buffer
	for i := 0; i < this.Length; i++ {
		result.WriteRune(this.Buffer[i])
	}
	return result.String()
}

func (this *Buffer) CurrentWordTop() (wordTop int) {
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

func (this *Buffer) CurrentWord() (string, int) {
	var buffer bytes.Buffer
	start := this.CurrentWordTop()
	for i := start; i < this.Cursor; i++ {
		if this.Buffer[i] != '"' {
			buffer.WriteRune(this.Buffer[i])
		}
	}
	return buffer.String(), start
}
