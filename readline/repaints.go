package readline

func (this *Buffer) InsertAndRepaint(str string) {
	this.ReplaceAndRepaint(this.Cursor, str)
}

func (this *Buffer) ReplaceAndRepaint(pos int, str string) {
	// Cursor rewind
	_, left, _ := this.view3()
	this.backspace(left.Width())

	// Replace Buffer
	this.Delete(pos, this.Cursor-pos)

	// Define ViewStart , Cursor
	this.Cursor = pos + this.InsertString(pos, str)
	this.ResetViewStart()

	// Repaint
	view, _, right := this.view3()
	this.puts(view)

	// Move to cursor position
	this.Eraseline()
	this.backspace(right.Width())

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
