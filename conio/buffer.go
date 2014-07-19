package conio

import "bufio"
import "bytes"
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
	Prompt    func() int
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

func (this *ReadLineBuffer) InsertAndRepaint(str string) {
	this.ReplaceAndRepaint(this.Cursor, str)
}

func (this *ReadLineBuffer) ReplaceAndRepaint(pos int, str string) {
	n := this.Cursor - pos
	if n > 0 {
		this.Delete(pos, n)
	}
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

func (this *ReadLineBuffer) RepaintAll() {
	if this.Prompt != nil {
		this.Prompt()
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
