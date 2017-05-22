package mains

import (
	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/shell"
)

type ShellVariable struct {
}

func (this *ShellVariable) Lookup(name string) string {
	f, ok := shell.PercentFunc[name]
	if ok {
		return f()
	} else {
		return ""
	}
}

func (this *ShellVariable) EachKey(f func(string)) {
	for name, _ := range shell.PercentFunc {
		f(name)
	}
}

func init() {
	completion.PercentVariables =
		append(completion.PercentVariables, new(ShellVariable))
}
