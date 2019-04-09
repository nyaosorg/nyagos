package readline

import (
	"fmt"
	"io"

	"github.com/mattn/go-runewidth"
)

var SurrogatePairOk = false

func (this *Buffer) putRune(ch rune) {
	if ch < ' ' {
		this.Out.WriteByte('^')
		this.Out.WriteByte(byte('A' + (ch - 1)))
	} else if (ch >= 0x10000 && !SurrogatePairOk) || runewidth.RuneWidth(ch) == 0 {
		fmt.Fprintf(this.Out, "<%X>", ch)
	} else {
		this.Out.WriteRune(ch)
	}
}

func (this *Buffer) putRunes(ch rune, n width_t) {
	if n <= 0 {
		return
	}
	this.putRune(ch)
	for i := width_t(1); i < n; i++ {
		this.Out.WriteRune(ch)
	}
}

func (this *Buffer) backspace(n width_t) {
	if n > 1 {
		fmt.Fprintf(this.Out, "\x1B[%dD", n)
	} else if n == 1 {
		this.Out.WriteByte('\b')
	}
}

func (this *Buffer) Eraseline() {
	io.WriteString(this.Out, "\x1B[0K")
}

type runes_t []rune

func (this *Buffer) puts(s []rune) runes_t {
	for _, ch := range s {
		this.putRune(ch)
	}
	return runes_t(s)
}

func (s runes_t) Width() (w width_t) {
	for _, ch := range s {
		w += GetCharWidth(ch)
	}
	return
}
