package readline

import (
	"fmt"

	"../conio"
)

type THistory struct {
	body []string
}

func (this *THistory) Len() int {
	return len(this.body)
}

func (this *THistory) At(n int) string {
	for n < 0 {
		n += len(this.body)
	}
	return this.body[n%len(this.body)]
}

func (this *THistory) Push(line string) {
	this.body = append(this.body, line)
}

func (this *THistory) Replace(line string) {
	if len(this.body) >= 1 {
		this.body[len(this.body)-1] = line
	} else {
		this.body = []string{line}
	}
}

type LineEditor struct {
	History THistory
	Pointer int
	Prompt  func(*LineEditor) (int, error)
	Tag     interface{}
}

func NewLineEditor() *LineEditor {
	return &LineEditor{
		Pointer: 0,
		Prompt:  func(this *LineEditor) (int, error) { fmt.Print("\n> "); return 2, nil },
	}
}

func (this *LineEditor) SetPromptStr(prompt string) {
	width := conio.GetStringWidth(prompt)
	this.Prompt = func(*LineEditor) (int, error) {
		fmt.Print(prompt)
		return width, nil
	}
}

var DefaultEditor *LineEditor

func KeyFuncHistoryUp(this *Buffer) Result {
	if this.Session.Pointer <= 0 {
		this.Session.Pointer = this.Session.History.Len()
	}
	this.Session.Pointer -= 1
	KeyFuncClear(this)
	if this.Session.Pointer >= 0 {
		this.InsertString(0, this.Session.History.At(this.Session.Pointer))
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(this)
	}
	return CONTINUE
}

func KeyFuncHistoryDown(this *Buffer) Result {
	this.Session.Pointer += 1
	if this.Session.Pointer > this.Session.History.Len() {
		this.Session.Pointer = 0
	}
	KeyFuncClear(this)
	if this.Session.Pointer < this.Session.History.Len() {
		this.InsertString(0, this.Session.History.At(this.Session.Pointer))
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(this)
	}
	return CONTINUE
}
