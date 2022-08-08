package frame

import (
	"github.com/nyaosorg/nyagos/internal/completion"
	"github.com/nyaosorg/nyagos/internal/shell"
)

type _ShellVariable struct {
}

func (*_ShellVariable) Lookup(name string) string {
	f, ok := shell.PercentFunc[name]
	if ok {
		return f()
	}
	return ""
}

func (*_ShellVariable) EachKey(f func(string)) {
	for name := range shell.PercentFunc {
		f(name)
	}
}

func init() {
	completion.PercentVariables =
		append(completion.PercentVariables, new(_ShellVariable))
}
