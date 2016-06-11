package readline

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"

	. "../conio"
)

func KeyFuncPass(this *Buffer) Result {
	return CONTINUE
}

func KeyFuncEnter(this *Buffer) Result { // Ctrl-M
	return ENTER
}

func KeyFuncIntr(this *Buffer) Result { // Ctrl-C
	this.Length = 0
	this.Cursor = 0
	this.ViewStart = 0
	return ENTER
}

func KeyFuncHead(this *Buffer) Result { // Ctrl-A
	Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
	this.Cursor = 0
	this.ViewStart = 0
	this.Repaint(0, 1)
	return CONTINUE
}

func KeyFuncBackword(this *Buffer) Result { // Ctrl-B
	if this.Cursor <= 0 {
		return CONTINUE
	}
	this.Cursor--
	if this.Cursor < this.ViewStart {
		this.ViewStart--
		this.Repaint(this.Cursor, 1)
	} else {
		Backspace(GetCharWidth(this.Buffer[this.Cursor]))
	}
	return CONTINUE
}

func KeyFuncTail(this *Buffer) Result { // Ctrl-E
	allength := this.GetWidthBetween(this.ViewStart, this.Length)
	if allength < this.ViewWidth {
		for ; this.Cursor < this.Length; this.Cursor++ {
			PutRune(this.Buffer[this.Cursor])
		}
	} else {
		PutRune('\a')
		Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
		this.ViewStart = this.Length - 1
		w := GetCharWidth(this.Buffer[this.ViewStart])
		for {
			if this.ViewStart <= 0 {
				break
			}
			w_ := w + GetCharWidth(this.Buffer[this.ViewStart-1])
			if w_ >= this.ViewWidth {
				break
			}
			w = w_
			this.ViewStart--
		}
		for this.Cursor = this.ViewStart; this.Cursor < this.Length; this.Cursor++ {
			PutRune(this.Buffer[this.Cursor])
		}
	}
	return CONTINUE
}

func KeyFuncForward(this *Buffer) Result { // Ctrl-F
	if this.Cursor >= this.Length {
		return CONTINUE
	}
	w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
	if w < this.ViewWidth {
		// No Scroll
		PutRune(this.Buffer[this.Cursor])
	} else {
		// Right Scroll
		Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
		if GetCharWidth(this.Buffer[this.Cursor]) > GetCharWidth(this.Buffer[this.ViewStart]) {
			this.ViewStart++
		}
		this.ViewStart++
		for i := this.ViewStart; i <= this.Cursor; i++ {
			PutRune(this.Buffer[i])
		}
		PutRune(' ')
		Backspace(1)
	}
	this.Cursor++
	return CONTINUE
}

func KeyFuncBackSpace(this *Buffer) Result { // Backspace
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

func KeyFuncDelete(this *Buffer) Result { // Del
	delw := this.Delete(this.Cursor, 1)
	this.Repaint(this.Cursor, delw)
	return CONTINUE
}

func KeyFuncDeleteOrAbort(this *Buffer) Result { // Ctrl-D
	if this.Length > 0 {
		return KeyFuncDelete(this)
	} else {
		return ABORT
	}
}

func KeyFuncInsertSelf(this *Buffer) Result {
	ch := this.Unicode
	if ch < 0x20 {
		return CONTINUE
	}
	this.Insert(this.Cursor, []rune{ch})

	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	w1 := GetCharWidth(ch)
	if w+w1 >= this.ViewWidth {
		// scroll left
		Backspace(w)
		if GetCharWidth(this.Buffer[this.ViewStart]) < w1 {
			this.ViewStart++
		}
		this.ViewStart++
		for i := this.ViewStart; i <= this.Cursor; i++ {
			PutRune(this.Buffer[i])
		}
		PutRune(' ')
		Backspace(1)
	} else {
		this.Repaint(this.Cursor, -w1)
	}
	this.Cursor++
	return CONTINUE
}

func KeyFuncInsertReport(this *Buffer) Result {
	this.InsertAndRepaint(fmt.Sprintf("[%X]", this.Unicode))
	return CONTINUE
}

func KeyFuncClearAfter(this *Buffer) Result {
	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	i := this.Cursor
	bs := 0

	var killbuf bytes.Buffer
	for j := this.Cursor; j < this.Length; j++ {
		killbuf.WriteRune(this.Buffer[j])
	}
	clipboard.WriteAll(killbuf.String())

	for i < this.Length && w < this.ViewWidth {
		w1 := GetCharWidth(this.Buffer[i])
		PutRunes(' ', w1)
		i++
		w += w1
		bs += w1
	}
	Backspace(bs)
	this.Length = this.Cursor
	return CONTINUE
}

func KeyFuncClear(this *Buffer) Result {
	width := this.GetWidthBetween(this.ViewStart, this.Cursor)
	Backspace(width)
	PutRunes(' ', this.ViewWidth)
	Backspace(this.ViewWidth)
	this.Length = 0
	this.Cursor = 0
	this.ViewStart = 0
	return CONTINUE
}

func KeyFuncWordRubout(this *Buffer) Result {
	org_cursor := this.Cursor
	for this.Cursor > 0 && unicode.IsSpace(this.Buffer[this.Cursor-1]) {
		this.Cursor--
	}
	i := this.CurrentWordTop()
	var killbuf bytes.Buffer
	for j := i; j < org_cursor; j++ {
		killbuf.WriteRune(this.Buffer[j])
	}
	clipboard.WriteAll(killbuf.String())
	keta := this.Delete(i, org_cursor-i)
	if i >= this.ViewStart {
		Backspace(keta)
	} else {
		Backspace(this.GetWidthBetween(this.ViewStart, org_cursor))
	}
	this.Cursor = i
	this.Repaint(i, keta)
	return CONTINUE
}

func KeyFuncClearBefore(this *Buffer) Result {
	keta := this.GetWidthBetween(this.ViewStart, this.Cursor)
	var killbuf bytes.Buffer
	for i := 0; i < this.Cursor; i++ {
		killbuf.WriteRune(this.Buffer[i])
	}
	clipboard.WriteAll(killbuf.String())
	this.Delete(0, this.Cursor)
	Backspace(keta)
	this.Cursor = 0
	this.ViewStart = 0
	this.Repaint(0, keta)
	return CONTINUE
}

func KeyFuncCLS(this *Buffer) Result {
	Cls()
	this.RepaintAll()
	return CONTINUE
}

func KeyFuncPaste(this *Buffer) Result {
	text, err := clipboard.ReadAll()
	if err == nil {
		this.InsertAndRepaint(
			strings.Replace(
				strings.Replace(
					strings.Replace(text, "\n", " ", -1),
					"\r", "", -1),
				"\t", " ", -1))
	}
	return CONTINUE
}

func maxInt(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func KeyFuncSwapChar(this *Buffer) Result {
	if this.Length == this.Cursor {
		if this.Cursor < 2 {
			return CONTINUE
		}
		this.Buffer[this.Cursor-2], this.Buffer[this.Cursor-1] = this.Buffer[this.Cursor-1], this.Buffer[this.Cursor-2]

		redrawStart := maxInt(this.Cursor-2, this.ViewStart)
		Backspace(this.GetWidthBetween(redrawStart, this.Cursor))
		for i := redrawStart; i < this.Cursor; i++ {
			PutRune(this.Buffer[i])
		}
	} else {
		if this.Cursor < 1 {
			return CONTINUE
		}

		w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
		this.Buffer[this.Cursor-1], this.Buffer[this.Cursor] = this.Buffer[this.Cursor], this.Buffer[this.Cursor-1]
		if w >= this.ViewWidth {
			// cursor move right and scroll
			w_1 := w - GetCharWidth(this.Buffer[this.Cursor])
			Backspace(w_1)
			this.ViewStart++
			for i := this.ViewStart; i <= this.Cursor; i++ {
				PutRune(this.Buffer[i])
			}
		} else {
			// no neccesary to scroll
			redrawStart := maxInt(this.Cursor-1, this.ViewStart)
			Backspace(this.GetWidthBetween(redrawStart, this.Cursor))
			for i := redrawStart; i <= this.Cursor; i++ {
				PutRune(this.Buffer[i])
			}
		}
		this.Cursor++
	}
	return CONTINUE
}
