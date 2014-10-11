package conio

import "bytes"
import "fmt"
import "strings"

import "github.com/atotto/clipboard"

func KeyFuncPass(this *ReadLineBuffer) KeyFuncResult {
	return CONTINUE
}

func KeyFuncEnter(this *ReadLineBuffer) KeyFuncResult { // Ctrl-M
	return ENTER
}

func KeyFuncIntr(this *ReadLineBuffer) KeyFuncResult { // Ctrl-C
	this.Length = 0
	this.Cursor = 0
	this.Buffer = []rune{}
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
		Backspace(GetCharWidth(this.Buffer[this.Cursor]))
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
		if GetCharWidth(this.Buffer[this.Cursor]) > GetCharWidth(this.Buffer[this.ViewStart]) {
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
	w1 := GetCharWidth(ch)
	if w+w1 >= this.ViewWidth {
		// scroll left
		Backspace(w)
		if GetCharWidth(this.Buffer[this.ViewStart]) < w1 {
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
	this.InsertAndRepaint(fmt.Sprintf("[%X]", this.Unicode))
	return CONTINUE
}

func KeyFuncClearAfter(this *ReadLineBuffer) KeyFuncResult {
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

func KeyFuncCLS(this *ReadLineBuffer) KeyFuncResult {
	Cls()
	this.RepaintAll()
	return CONTINUE
}

func KeyFuncPaste(this *ReadLineBuffer) KeyFuncResult {
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
