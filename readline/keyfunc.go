package readline

import (
	"context"
	"io"
	"strings"
	"unicode"

	"github.com/atotto/clipboard"
)

func KeyFuncEnter(ctx context.Context, this *Buffer) Result { // Ctrl-M
	return ENTER
}

func KeyFuncIntr(ctx context.Context, this *Buffer) Result { // Ctrl-C
	this.Length = 0
	this.Cursor = 0
	this.ViewStart = 0
	return INTR
}

func KeyFuncHead(ctx context.Context, this *Buffer) Result { // Ctrl-A
	this.Backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))
	this.Cursor = 0
	this.ViewStart = 0
	this.Repaint(0, 1)
	return CONTINUE
}

func KeyFuncBackword(ctx context.Context, this *Buffer) Result { // Ctrl-B
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

func KeyFuncTail(ctx context.Context, this *Buffer) Result { // Ctrl-E
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

func KeyFuncForward(ctx context.Context, this *Buffer) Result { // Ctrl-F
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

func KeyFuncBackSpace(ctx context.Context, this *Buffer) Result { // Backspace
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

func KeyFuncDelete(ctx context.Context, this *Buffer) Result { // Del
	delw := this.Delete(this.Cursor, 1)
	this.Repaint(this.Cursor, delw)
	return CONTINUE
}

func KeyFuncDeleteOrAbort(ctx context.Context, this *Buffer) Result { // Ctrl-D
	if this.Length > 0 {
		return KeyFuncDelete(ctx, this)
	} else {
		return ABORT
	}
}

func KeyFuncInsertSelf(ctx context.Context, this *Buffer, keys string) Result {
	this.Insert(this.Cursor, []rune(keys))

	w := this.GetWidthBetween(this.ViewStart, this.Cursor)
	w1 := GetStringWidth(keys)
	if w+w1 >= this.ViewWidth() {
		// scroll left
		this.Backspace(w)
		this.Cursor++
		this.ResetViewStart()
		for i := this.ViewStart; i < this.Cursor; i++ {
			this.PutRune(this.Buffer[i])
		}
		this.Eraseline()
	} else {
		this.Repaint(this.Cursor, -w1)
		this.Cursor++
	}
	return CONTINUE
}

func KeyFuncClearAfter(ctx context.Context, this *Buffer) Result {
	clipboard.WriteAll(this.SubString(this.Cursor, this.Length))

	this.Eraseline()
	this.Length = this.Cursor
	return CONTINUE
}

func KeyFuncClear(ctx context.Context, this *Buffer) Result {
	width := this.GetWidthBetween(this.ViewStart, this.Cursor)
	this.Backspace(width)
	this.Eraseline()
	this.Length = 0
	this.Cursor = 0
	this.ViewStart = 0
	return CONTINUE
}

func KeyFuncWordRubout(ctx context.Context, this *Buffer) Result {
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

func KeyFuncClearBefore(ctx context.Context, this *Buffer) Result {
	keta := this.GetWidthBetween(this.ViewStart, this.Cursor)
	clipboard.WriteAll(this.SubString(0, this.Cursor))
	this.Delete(0, this.Cursor)
	this.Backspace(keta)
	this.Cursor = 0
	this.ViewStart = 0
	this.Repaint(0, keta)
	return CONTINUE
}

func KeyFuncCLS(ctx context.Context, this *Buffer) Result {
	io.WriteString(this.Out, "\x1B[1;1H\x1B[2J")
	this.RepaintAll()
	return CONTINUE
}

func KeyFuncRepaintOnNewline(ctx context.Context, this *Buffer) Result {
	this.Out.WriteByte('\n')
	this.RepaintAll()
	return CONTINUE
}

func KeyFuncQuotedInsert(ctx context.Context, this *Buffer) Result {
	io.WriteString(this.Out, CURSOR_ON)
	defer io.WriteString(this.Out, CURSOR_OFF)

	this.Out.Flush()
	if key, err := getKey(this.TTY); err == nil {
		return KeyFuncInsertSelf(ctx, this, key)
	} else {
		return CONTINUE
	}
}

func KeyFuncPaste(ctx context.Context, this *Buffer) Result {
	text, err := clipboard.ReadAll()
	if err != nil {
		return CONTINUE
	}
	this.InsertAndRepaint(text)
	return CONTINUE
}

func KeyFuncPasteQuote(ctx context.Context, this *Buffer) Result {
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

func KeyFuncSwapChar(ctx context.Context, this *Buffer) Result {
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
