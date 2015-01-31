package conio

import "fmt"

type HistoryLine struct {
	Line string
	Word []string
}

func NewHistoryLine(line string) *HistoryLine {
	return &HistoryLine{Line: line, Word: SplitQ(line)}
}

func (this *HistoryLine) At(n int) string {
	for n < 0 {
		n += len(this.Word)
	}
	return this.Word[n%len(this.Word)]
}

type LineEditor struct {
	Histories []*HistoryLine
	Pointer   int
	Prompt    func() int
}

func NewLineEditor() *LineEditor {
	return &LineEditor{
		Histories: make([]*HistoryLine, 0),
		Pointer:   0,
		Prompt:    func() int { fmt.Print("\n> "); return 2 },
	}
}

func (this *LineEditor) GetHistoryAt(n int) *HistoryLine {
	if n < 0 {
		n = len(this.Histories) + n
	}
	if n >= len(this.Histories) {
		return &HistoryLine{Line: "", Word: []string{}}
	} else {
		return this.Histories[n]
	}
}

func (this *LineEditor) HistoryLen() int {
	return len(this.Histories)
}

func (this *LineEditor) LastHistory() *HistoryLine {
	if len(this.Histories) <= 0 {
		return nil
	} else {
		return this.Histories[len(this.Histories)-1]
	}
}

func (this *LineEditor) HistoryPush(input string) {
	this.Histories = append(this.Histories, NewHistoryLine(input))
	this.HistoryResetPointer()
}

func (this *LineEditor) HistoryResetPointer() {
	this.Pointer = len(this.Histories)
}

func (this *LineEditor) SetPromptStr(prompt string) {
	width := GetStringWidth(prompt)
	this.Prompt = func() int {
		fmt.Print(prompt)
		return width
	}
}

var DefaultEditor = NewLineEditor()

func KeyFuncHistoryUp(this *Buffer) Result {
	if this.Session.Pointer <= 0 {
		this.Session.Pointer = this.Session.HistoryLen()
	}
	this.Session.Pointer -= 1
	KeyFuncClear(this)
	if this.Session.Pointer >= 0 {
		this.InsertString(0, this.Session.Histories[this.Session.Pointer].Line)
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(this)
	}
	return CONTINUE
}

func KeyFuncHistoryDown(this *Buffer) Result {
	this.Session.Pointer += 1
	if this.Session.Pointer > this.Session.HistoryLen() {
		this.Session.Pointer = 0
	}
	KeyFuncClear(this)
	if this.Session.Pointer < this.Session.HistoryLen() {
		this.InsertString(0, this.Session.Histories[this.Session.Pointer].Line)
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(this)
	}
	return CONTINUE
}
