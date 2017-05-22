package completion

import (
	"github.com/zetamatta/nyagos/readline"
)

var command_listupper = []func() []Element{
	func() []Element { return listUpAllExecutableOnEnv("PATH") },
	func() []Element { return listUpAllExecutableOnEnv("NYAGOSPATH") },
}

func AppendCommandLister(f func() []Element) {
	command_listupper = append(command_listupper, f)
}

var HookToList = []func(*readline.Buffer, *List) (*List, error){}

func init() {
	f := readline.KeyGoFuncT{Func: KeyFuncCompletion, Name: "COMPLETE"}
	err := readline.BindKeyFunc(readline.K_CTRL_I, &f)
	if err != nil {
		panic(err.Error())
	}
}
