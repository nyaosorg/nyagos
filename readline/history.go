package readline

import (
	"fmt"

	"../conio"
)

type IHistory interface {
	Len() int
	At(int) string
	Push(string)
	Replace(string)
}

type LineEditor struct {
	History IHistory
	Pointer int
	Prompt  func(*LineEditor) (int, error)
	Tag     interface{}
}

func NewLineEditor(history IHistory) *LineEditor {
	return &LineEditor{
		History: history,
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
