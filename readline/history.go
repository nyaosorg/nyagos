package readline

type IHistory interface {
	Len() int
	At(int) string
}

type Editor struct {
	History IHistory
	Prompt  func() (int, error)
}

func KeyFuncHistoryUp(this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	if this.HistoryPointer <= 0 {
		this.HistoryPointer = this.History.Len()
	}
	this.HistoryPointer -= 1
	KeyFuncClear(this)
	if this.HistoryPointer >= 0 {
		this.InsertString(0, this.History.At(this.HistoryPointer))
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(this)
	}
	return CONTINUE
}

func KeyFuncHistoryDown(this *Buffer) Result {
	if this.History.Len() <= 0 {
		return CONTINUE
	}
	this.HistoryPointer += 1
	if this.HistoryPointer > this.History.Len() {
		this.HistoryPointer = 0
	}
	KeyFuncClear(this)
	if this.HistoryPointer < this.History.Len() {
		this.InsertString(0, this.History.At(this.HistoryPointer))
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(this)
	}
	return CONTINUE
}
