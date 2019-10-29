package readline

import (
	"bufio"
	"context"
	"io"
)

type IHistory interface {
	Len() int
	At(int) string
}

type Editor struct {
	History  IHistory
	Writer   io.Writer
	Out      *bufio.Writer
	Prompt   func() (int, error)
	Default  string
	Cursor   int
	LineFeed func(Result)
}

func keyFuncHistoryUp(ctx context.Context, this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	if this.HistoryPointer <= 0 {
		this.HistoryPointer = this.History.Len()
	}
	this.HistoryPointer -= 1
	keyFuncClear(ctx, this)
	if this.HistoryPointer >= 0 {
		this.InsertString(0, this.History.At(this.HistoryPointer))
		this.ViewStart = 0
		this.Cursor = 0
		keyFuncTail(ctx, this)
	}
	return CONTINUE
}

func keyFuncHistoryDown(ctx context.Context, this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	if this.HistoryPointer+1 > this.History.Len() {
		return CONTINUE
	}
	this.HistoryPointer += 1
	keyFuncClear(ctx, this)
	if this.HistoryPointer < this.History.Len() {
		this.InsertString(0, this.History.At(this.HistoryPointer))
		this.ViewStart = 0
		this.Cursor = 0
		keyFuncTail(ctx, this)
	}
	return CONTINUE
}
