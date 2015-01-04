package conio

var Histories = make([]string, 0)
var HistoryPointer = 0

func GetHistory(n int) string {
	if n < 0 {
		n = len(Histories) + n
	}
	if n >= len(Histories) {
		return ""
	} else {
		return Histories[n]
	}
}

func HistoryLen() int {
	return len(Histories)
}

func LastHistory() string {
	if len(Histories) <= 0 {
		return ""
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
		this.InsertString(0, Histories[HistoryPointer])
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
		this.InsertString(0, Histories[HistoryPointer])
		this.ViewStart = 0
		this.Cursor = 0
		KeyFuncTail(this)
	}
	return CONTINUE
}

func HistoryPush(input string) {
	Histories = append(Histories, input)
	HistoryResetPointer()
}

func HistoryResetPointer() {
	HistoryPointer = len(Histories)
}
