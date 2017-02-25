package completion

import (
	"../readline"
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
	f := readline.KeyGoFuncT{F: KeyFuncCompletion}
	err := readline.BindKeySymbolFunc(readline.K_CTRL_I, "COMPLETE", &f)
	if err != nil {
		panic(err.Error())
	}
}
