package readline

import "bytes"
import "fmt"
import "strings"

import "github.com/atotto/clipboard"
import ".."

func KeyFuncPass(this *Buffer) KeyFuncResult {
	return CONTINUE
}

func KeyFuncEnter(this *Buffer) KeyFuncResult { // Ctrl-M
	return ENTER
}

func KeyFuncIntr(this *Buffer) KeyFuncResult { // Ctrl-C
	this.Length = 0
	this.Cursor = 0
	this.Buffer = []rune{}
	return ENTER
}

func KeyFuncHead(this *Buffer) KeyFuncResult { // Ctrl-A
	Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
	this.Cursor = 0
	this.ViewStart = 0
	this.Repaint(0, 1)
	return CONTINUE
}

func KeyFuncBackword(this *Buffer) KeyFuncResult { // Ctrl-B
	if this.Cursor <= 0 {
		return CONTINUE
	}
	this.Cursor--
	if this.Cursor < this.ViewStart {
		this.ViewStart--
		this.Repaint(this.Cursor, 1)
	} else {
		Backspace(conio.GetCharWidth(this.Buffer[this.Cursor]))
	}
	return CONTINUE
}

func KeyFuncTail(this *Buffer) KeyFuncResult { // Ctrl-E
	allength := this.GetWidthBetween(this.ViewStart, this.Length)
	if allength < this.ViewWidth {
		for ; this.Cursor < this.Length; this.Cursor++ {
			PutRep(this.Buffer[this.Cursor], 1)
		}
	} else {
		PutRep('\a', 1)
		Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
		this.ViewStart = this.Length - 1
		w := conio.GetCharWidth(this.Buffer[this.ViewStart])
		for {
			if this.ViewStart <= 0 {
				break
			}
			w_ := w + conio.GetCharWidth(this.Buffer[this.ViewStart-1])
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

func KeyFuncForward(this *Buffer) KeyFuncResult { // Ctrl-F
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
		if conio.GetCharWidth(this.Buffer[this.Cursor]) > conio.GetCharWidth(this.Buffer[this.ViewStart]) {
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

func KeyFuncBackSpace(this *Buffer) KeyFuncResult { // Backspace
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

func KeyFuncDelete(this *Buffer) KeyFuncResult { // Del
	delw := this.Delete(this.Cursor, 1)
	this.Repaint(this.Cursor, delw)
	return CONTINUE
}

func KeyFuncDeleteOrAbort(this *Buffer) KeyFuncResult { // Ctrl-D
	if this.Length > 0 {
		return KeyFuncDelete(this)
	} else {
		return ABORT
	}
}

func KeyFuncInsertSelf(this *Buffer) KeyFuncResult {
	ch := this.Unicode
	if ch < 0x20 || !this.Insert(this.Cursor, []rune{ch}) {
		return CONTINUE
	}
	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	w1 := conio.GetCharWidth(ch)
	if w+w1 >= this.ViewWidth {
		// scroll left
		Backspace(w)
		if conio.GetCharWidth(this.Buffer[this.ViewStart]) < w1 {
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

func KeyFuncInsertReport(this *Buffer) KeyFuncResult {
	this.InsertAndRepaint(fmt.Sprintf("[%X]", this.Unicode))
	return CONTINUE
}

func KeyFuncClearAfter(this *Buffer) KeyFuncResult {
	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	i := this.Cursor
	bs := 0

	var killbuf bytes.Buffer
	for j := this.Cursor; j < this.Length; j++ {
		killbuf.WriteRune(this.Buffer[j])
	}
	clipboard.WriteAll(killbuf.String())

	for i < this.Length && w < this.ViewWidth {
		w1 := conio.GetCharWidth(this.Buffer[i])
		PutRep(' ', w1)
		i++
		w += w1
		bs += w1
	}
	Backspace(bs)
	this.Length = this.Cursor
	return CONTINUE
}

func KeyFuncClear(this *Buffer) KeyFuncResult {
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

func KeyFuncClearBefore(this *Buffer) KeyFuncResult {
	keta := this.Delete(0, this.Cursor)
	Backspace(keta)
	this.Cursor = 0
	this.Repaint(0, keta)
	return CONTINUE
}

func KeyFuncCLS(this *Buffer) KeyFuncResult {
	conio.Cls()
	this.RepaintAll()
	return CONTINUE
}

func KeyFuncPaste(this *Buffer) KeyFuncResult {
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
