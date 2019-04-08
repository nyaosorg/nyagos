package readline

import (
	"context"
	"io"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"
)

func keyFuncEnter(ctx context.Context, this *Buffer) Result { // Ctrl-M
	return ENTER
}

func keyFuncIntr(ctx context.Context, this *Buffer) Result { // Ctrl-C
	this.Length = 0
	this.Cursor = 0
	this.ViewStart = 0
	return INTR
}

func keyFuncHead(ctx context.Context, this *Buffer) Result { // Ctrl-A
	this.Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
	this.Cursor = 0
	this.ViewStart = 0
	this.Repaint(0, 1)
	return CONTINUE
}

func keyFuncBackward(ctx context.Context, this *Buffer) Result { // Ctrl-B
	if this.Cursor <= 0 {
		return CONTINUE
	}
	this.Cursor--
	if this.Cursor < this.ViewStart {
		this.ViewStart--
		this.Repaint(this.Cursor, 1)
	} else {
		this.Backspace(GetCharWidth(this.Buffer[this.Cursor]))
	}
	return CONTINUE
}

func keyFuncTail(ctx context.Context, this *Buffer) Result { // Ctrl-E
	allength := this.GetWidthBetween(this.ViewStart, this.Length)
	if allength < this.ViewWidth() {
		for ; this.Cursor < this.Length; this.Cursor++ {
			this.PutRune(this.Buffer[this.Cursor])
		}
	} else {
		io.WriteString(this.Out, "\a")
		this.Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
		this.ViewStart = this.Length - 1
		w := GetCharWidth(this.Buffer[this.ViewStart])
		for {
			if this.ViewStart <= 0 {
				break
			}
			w_ := w + GetCharWidth(this.Buffer[this.ViewStart-1])
			if w_ >= this.ViewWidth() {
				break
			}
			w = w_
			this.ViewStart--
		}
		for this.Cursor = this.ViewStart; this.Cursor < this.Length; this.Cursor++ {
			this.PutRune(this.Buffer[this.Cursor])
		}
	}
	return CONTINUE
}

func keyFuncForward(ctx context.Context, this *Buffer) Result { // Ctrl-F
	if this.Cursor >= this.Length {
		return CONTINUE
	}
	w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
	if w < this.ViewWidth() {
		// No Scroll
		this.PutRune(this.Buffer[this.Cursor])
	} else {
		// Right Scroll
		this.Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
		if GetCharWidth(this.Buffer[this.Cursor]) > GetCharWidth(this.Buffer[this.ViewStart]) {
			this.ViewStart++
		}
		this.ViewStart++
		for i := this.ViewStart; i <= this.Cursor; i++ {
			this.PutRune(this.Buffer[i])
		}
		this.Eraseline()
	}
	this.Cursor++
	return CONTINUE
}

func keyFuncBackSpace(ctx context.Context, this *Buffer) Result { // Backspace
	if this.Cursor > 0 {
		this.Cursor--
		delw := this.Delete(this.Cursor, 1)
		if this.Cursor >= this.ViewStart {
			this.Backspace(delw)
		} else {
			this.ViewStart = this.Cursor
		}
		this.Repaint(this.Cursor, delw)
	}
	return CONTINUE
}

func keyFuncDelete(ctx context.Context, this *Buffer) Result { // Del
	delw := this.Delete(this.Cursor, 1)
	this.Repaint(this.Cursor, delw)
	return CONTINUE
}

func keyFuncDeleteOrAbort(ctx context.Context, this *Buffer) Result { // Ctrl-D
	if this.Length > 0 {
		return keyFuncDelete(ctx, this)
	} else {
		return ABORT
	}
}

func keyFuncInsertSelf(ctx context.Context, this *Buffer, keys string) Result {
	if len(keys) == 2 && keys[0] == '\x1B' { // for AltGr-shift
		keys = keys[1:]
	}
	this.Insert(this.Cursor, []rune(keys))

	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	w1 := GetStringWidth(keys)
	if w+w1 >= this.ViewWidth() {
		// scroll left
		this.Backspace(w)
		this.Cursor += len([]rune(keys))
		this.ResetViewStart()
		for i := this.ViewStart; i < this.Cursor; i++ {
			this.PutRune(this.Buffer[i])
		}
		this.Eraseline()
	} else {
		this.Repaint(this.Cursor, -w1)
		this.Cursor += len([]rune(keys))
	}
	return CONTINUE
}

func keyFuncClearAfter(ctx context.Context, this *Buffer) Result {
	clipboard.WriteAll(this.SubString(this.Cursor, this.Length))

	this.Eraseline()
	this.Length = this.Cursor
	return CONTINUE
}

func keyFuncClear(ctx context.Context, this *Buffer) Result {
	width := this.GetWidthBetween(this.ViewStart, this.Cursor)
	this.Backspace(width)
	this.Eraseline()
	this.Length = 0
	this.Cursor = 0
	this.ViewStart = 0
	return CONTINUE
}

func keyFuncWordRubout(ctx context.Context, this *Buffer) Result {
	org_cursor := this.Cursor
	for this.Cursor > 0 && unicode.IsSpace(this.Buffer[this.Cursor-1]) {
		this.Cursor--
	}
	i := this.CurrentWordTop()
	clipboard.WriteAll(this.SubString(i, org_cursor))
	keta := this.Delete(i, org_cursor-i)
	if i >= this.ViewStart {
		this.Backspace(keta)
	} else {
		this.Backspace(this.GetWidthBetween(this.ViewStart, org_cursor))
	}
	this.Cursor = i
	this.Repaint(i, keta)
	return CONTINUE
}

func keyFuncClearBefore(ctx context.Context, this *Buffer) Result {
	keta := this.GetWidthBetween(this.ViewStart, this.Cursor)
	clipboard.WriteAll(this.SubString(0, this.Cursor))
	this.Delete(0, this.Cursor)
	this.Backspace(keta)
	this.Cursor = 0
	this.ViewStart = 0
	this.Repaint(0, keta)
	return CONTINUE
}

func keyFuncCLS(ctx context.Context, this *Buffer) Result {
	io.WriteString(this.Out, "\x1B[1;1H\x1B[2J")
	this.RepaintAll()
	return CONTINUE
}

func keyFuncRepaintOnNewline(ctx context.Context, this *Buffer) Result {
	this.Out.WriteByte('\n')
	this.RepaintAll()
	return CONTINUE
}

func keyFuncQuotedInsert(ctx context.Context, this *Buffer) Result {
	io.WriteString(this.Out, ansiCursorOn)
	defer io.WriteString(this.Out, ansiCursorOff)

	this.Out.Flush()
	if key, err := getKey(this.TTY); err == nil {
		return keyFuncInsertSelf(ctx, this, key)
	} else {
		return CONTINUE
	}
}

func keyFuncPaste(ctx context.Context, this *Buffer) Result {
	text, err := clipboard.ReadAll()
	if err != nil {
		return CONTINUE
	}
	this.InsertAndRepaint(text)
	return CONTINUE
}

func keyFuncPasteQuote(ctx context.Context, this *Buffer) Result {
	text, err := clipboard.ReadAll()
	if err != nil {
		return CONTINUE
	}
	if strings.IndexRune(text, ' ') >= 0 &&
		!strings.HasPrefix(text, `"`) {
		text = `"` + strings.Replace(text, `"`, `""`, -1) + `"`
	}
	this.InsertAndRepaint(text)
	return CONTINUE
}

func maxInt(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func keyFuncSwapChar(ctx context.Context, this *Buffer) Result {
	if this.Length == this.Cursor {
		if this.Cursor < 2 {
			return CONTINUE
		}
		this.Buffer[this.Cursor-2], this.Buffer[this.Cursor-1] = this.Buffer[this.Cursor-1], this.Buffer[this.Cursor-2]

		redrawStart := maxInt(this.Cursor-2, this.ViewStart)
		this.Backspace(this.GetWidthBetween(redrawStart, this.Cursor))
		for i := redrawStart; i < this.Cursor; i++ {
			this.PutRune(this.Buffer[i])
		}
	} else {
		if this.Cursor < 1 {
			return CONTINUE
		}

		w := this.GetWidthBetween(this.ViewStart, this.Cursor+1)
		this.Buffer[this.Cursor-1], this.Buffer[this.Cursor] = this.Buffer[this.Cursor], this.Buffer[this.Cursor-1]
		if w >= this.ViewWidth() {
			// cursor move right and scroll
			w_1 := w - GetCharWidth(this.Buffer[this.Cursor])
			this.Backspace(w_1)
			this.ViewStart++
			for i := this.ViewStart; i <= this.Cursor; i++ {
				this.PutRune(this.Buffer[i])
			}
		} else {
			// no necessary to scroll
			redrawStart := maxInt(this.Cursor-1, this.ViewStart)
			this.Backspace(this.GetWidthBetween(redrawStart, this.Cursor))
			for i := redrawStart; i <= this.Cursor; i++ {
				this.PutRune(this.Buffer[i])
			}
		}
		this.Cursor++
	}
	return CONTINUE
}

func keyFuncBackwardWord(ctx context.Context, this *Buffer) Result {
	newPos := this.Cursor
	for newPos > 0 && this.Buffer[newPos-1] == ' ' {
		newPos--
	}
	for newPos > 0 && this.Buffer[newPos-1] != ' ' {
		newPos--
	}
	if newPos >= this.ViewStart {
		w := this.GetWidthBetween(newPos, this.Cursor)
		this.Backspace(w)
		this.Cursor = newPos
	} else {
		w := this.GetWidthBetween(this.ViewStart, this.Cursor)
		this.Backspace(w)
		this.Cursor = newPos
		this.ViewStart = newPos
		this.Repaint(newPos, 0)
	}
	return CONTINUE
}

func keyFuncForwardWord(ctx context.Context, this *Buffer) Result {
	newPos := this.Cursor
	for newPos < this.Length && this.Buffer[newPos] != ' ' {
		newPos++
	}
	for newPos < this.Length && this.Buffer[newPos] == ' ' {
		newPos++
	}
	w := this.GetWidthBetween(this.ViewStart, newPos)
	if w < this.ViewWidth() {
		for this.Cursor < newPos {
			this.PutRune(this.Buffer[this.Cursor])
			this.Cursor++
		}
	} else {
		this.Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
		this.Cursor = newPos
		for w >= this.ViewWidth() {
			w -= GetCharWidth(this.Buffer[this.ViewStart])
			this.ViewStart++
		}
		for p := this.ViewStart; p < this.Cursor; p++ {
			this.PutRune(this.Buffer[p])
		}
		this.Eraseline()
	}
	return CONTINUE
}
