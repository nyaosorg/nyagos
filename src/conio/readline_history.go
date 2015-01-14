package conio

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

var Histories = make([]*HistoryLine, 0)
var HistoryPointer = 0

func GetHistory(n int) *HistoryLine {
	if n < 0 {
		n = len(Histories) + n
	}
	if n >= len(Histories) {
		return &HistoryLine{Line: "", Word: []string{}}
	} else {
		return Histories[n]
	}
}

func HistoryLen() int {
	return len(Histories)
}

func LastHistory() *HistoryLine {
	if len(Histories) <= 0 {
		return nil
	} else {
		return Histories[len(Histories)-1]
	}
}

func KeyFuncHistoryUp(this *Buffer) Result {
	if HistoryPointer <= 0 {
		HistoryPointer = len(Histories)
	}
	HistoryPointer -= 1
	KeyFuncClear(this)
	if HistoryPointer >= 0 {
		this.InsertString(0, Histories[HistoryPointer].Line)
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(this)
	}
	return CONTINUE
}

func KeyFuncHistoryDown(this *Buffer) Result {
	HistoryPointer += 1
	if HistoryPointer > len(Histories) {
		HistoryPointer = 0
	}
	KeyFuncClear(this)
	if HistoryPointer < len(Histories) {
		this.InsertString(0, Histories[HistoryPointer].Line)
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(this)
	}
	return CONTINUE
}

func HistoryPush(input string) {
	Histories = append(Histories, NewHistoryLine(input))
	HistoryResetPointer()
}

func HistoryResetPointer() {
	HistoryPointer = len(Histories)
}
