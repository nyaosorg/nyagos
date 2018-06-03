package readline

import (
	"bufio"
	"context"
)

type IHistory interface {
	Len() int
	At(int) string
}

type Editor struct {
	History IHistory
	Writer  *bufio.Writer
	Prompt  func() (int, error)
	Default string
	Cursor  int
}

func KeyFuncHistoryUp(ctx context.Context, this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	if this.HistoryPointer <= 0 {
		this.HistoryPointer = this.History.Len()
	}
	this.HistoryPointer -= 1
	KeyFuncClear(ctx, this)
	if this.HistoryPointer >= 0 {
		this.InsertString(0, this.History.At(this.HistoryPointer))
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(ctx, this)
	}
	return CONTINUE
}

func KeyFuncHistoryDown(ctx context.Context, this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	if this.HistoryPointer+1 > this.History.Len() {
		return CONTINUE
	}
	this.HistoryPointer += 1
	KeyFuncClear(ctx, this)
	if this.HistoryPointer < this.History.Len() {
		this.InsertString(0, this.History.At(this.HistoryPointer))
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(ctx, this)
	}
	return CONTINUE
}
