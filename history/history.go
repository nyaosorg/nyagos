package history

import "../conio"

var histories = make([]string, 0)
var pointor = 0

func KeyFuncHistoryUp(this *conio.ReadLineBuffer) conio.KeyFuncResult {
	if pointor <= 0 {
		pointor = len(histories)
	}
	pointor -= 1
	conio.KeyFuncClear(this)
	if pointor >= 0 {
		this.InsertString(0, histories[pointor])
		this.ViewStart = 0
		this.Cursor = 0
		conio.KeyFuncTail(this)
	}
	return conio.CONTINUE
}

func KeyFuncHistoryDown(this *conio.ReadLineBuffer) conio.KeyFuncResult {
	pointor += 1
	if pointor >= len(histories) {
		pointor = 0
	}
	conio.KeyFuncClear(this)
	if pointor < len(histories) {
		this.InsertString(0, histories[pointor])
		this.ViewStart = 0
		this.Cursor = 0
		conio.KeyFuncTail(this)
	}
	return conio.CONTINUE
}

func Push(input string) {
	histories = append(histories, input)
	pointor = len(histories)
}
