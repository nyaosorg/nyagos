package main

import (
	"../readline"
)

type THistory struct {
	Obj *readline.LineEditor
}

func (this *THistory) Len() int {
	return this.Obj.History.Len()
}

func (this *THistory) At(n int) string {
	return this.Obj.History.At(n)
}

func (this *THistory) Push(line string) {
	this.Obj.History.Push(line)
}

func (this *THistory) Replace(line string) {
	this.Obj.History.Replace(line)
}
