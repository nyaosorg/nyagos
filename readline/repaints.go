package readline

func (this *Buffer) InsertAndRepaint(str string) {
	this.ReplaceAndRepaint(this.Cursor, str)
}

func (this *Buffer) ReplaceAndRepaint(pos int, str string) {
	// Cursor rewind
	this.backspace(this.GetWidthBetween(this.ViewStart, this.Cursor))

	// Replace Buffer
	this.Delete(pos, this.Cursor-pos)

	// Define ViewStart , Cursor
	this.Cursor = pos + this.InsertString(pos, str)
	this.ResetViewStart()

	// Repaint
	w := width_t(0)
	for _, ch := range this.Buffer[this.ViewStart:this.Cursor] {
		this.putRune(ch)
		w += GetCharWidth(ch)
	}
	bs := width_t(0)
	for _, ch := range this.Buffer[this.Cursor:] {
		w1 := GetCharWidth(ch)
		if w+w1 >= this.ViewWidth() {
			break
		}
		this.putRune(ch)
		w += w1
		bs += w1
	}
	this.Eraseline()
	if bs > 0 {
		this.backspace(bs)
	}
}

// Repaint buffer[pos:] + " \b"*del but do not rewind cursor position
func (this *Buffer) Repaint(pos int, del width_t) {
	bs := width_t(0)
	vp := this.GetWidthBetween(this.ViewStart, pos)

	for _, ch := range this.Buffer[pos:] {
		w1 := GetCharWidth(ch)
		if vp+w1 >= this.ViewWidth() {
			break
		}
		this.putRune(ch)
		vp += w1
		bs += w1
	}
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
	for _, ch := range this.Buffer[this.ViewStart:this.Cursor] {
		this.putRune(ch)
	}
	this.Repaint(this.Cursor, 0)
}

func (this *Buffer) RepaintAll() {
	this.Out.Flush()
	this.TopColumn, _ = this.Prompt()
	this.RepaintAfterPrompt()
}
