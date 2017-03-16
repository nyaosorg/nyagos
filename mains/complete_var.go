package mains

import (
	"../completion"
	"../interpreter"
)

type ShellVariable struct {
}

func (this *ShellVariable) Lookup(name string) string {
	f, ok := interpreter.PercentFunc[name]
	if ok {
		return f()
	} else {
		return ""
	}
}

func (this *ShellVariable) EachKey(f func(string)) {
	for name, _ := range interpreter.PercentFunc {
		f(name)
	}
}

func init() {
	completion.PercentVariables =
		append(completion.PercentVariables, new(ShellVariable))
}
