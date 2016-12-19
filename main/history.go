package main

import (
	"../history"
	"../readline"
)

type THistory struct {
	Obj *readline.LineEditor
}

func (this *THistory) Len() int {
	return this.Obj.HistoryLen()
}

func (this *THistory) At(n int) string {
	return this.Obj.Histories[n].Line
}

func (this *THistory) Push(line string) {
	tmp := readline.NewHistoryLine(line)
	this.Obj.Histories = append(this.Obj.Histories, tmp)
}

func (this *THistory) Replace(line string) {
	tmp := readline.NewHistoryLine(line)
	this.Obj.Histories[len(this.Obj.Histories)-1] = tmp
}

func historyReplace(line string) (string, bool) {
	hisObj := THistory{readline.DefaultEditor}
	return history.Replace(&hisObj, line)
}
