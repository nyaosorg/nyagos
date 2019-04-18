package readline

func (this *Buffer) InsertAndRepaint(str string) {
	this.ReplaceAndRepaint(this.Cursor, str)
}

// GotoHead move screen-cursor to the top of the viewarea.
// It should be called before text is changed.
func (this *Buffer) GotoHead() {
	this.backspace(range_t(this.Buffer[this.ViewStart:this.Cursor]).Width())
}

// DrawFromHead draw all text in viewarea and
// move screen-cursor to the position where it should be.
func (this *Buffer) DrawFromHead() {
	// Repaint
	view, _, right := this.view3()
	this.puts(view)

	// Move to cursor position
	this.Eraseline()
	this.backspace(right.Width())
}

func (this *Buffer) ReplaceAndRepaint(pos int, str string) {
	this.GotoHead()

	// Replace Buffer
	this.Delete(pos, this.Cursor-pos)

	// Define ViewStart , Cursor
	this.Cursor = pos + this.InsertString(pos, str)

	this.joinUndo() // merge delete and insert

	this.ResetViewStart()

	this.DrawFromHead()
}

// Repaint buffer[pos:] + " \b"*del but do not rewind cursor position
func (this *Buffer) Repaint(pos int, del width_t) {
	vp := this.GetWidthBetween(this.ViewStart, pos)

	view := this.view()
	bs := this.puts(view[pos-this.ViewStart:]).Width()
	vp += bs

	this.Eraseline()
	if del > 0 {
		this.backspace(bs)
	} else {
		// for readline_keyfunc.go: KeyFuncInsertSelf()
		this.backspace(bs + del)
	}
}

func (this *Buffer) RepaintAfterPrompt() {
	this.ResetViewStart()
	this.puts(this.Buffer[this.ViewStart:this.Cursor])
	this.Repaint(this.Cursor, 0)
}

func (this *Buffer) RepaintAll() {
	this.Out.Flush()
	this.TopColumn, _ = this.Prompt()
	this.RepaintAfterPrompt()
}
