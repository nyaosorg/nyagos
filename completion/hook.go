package completion

import (
	"context"
	"github.com/zetamatta/nyagos/readline"
)

var commandListUpper = []func() []Element{
	func() []Element { return listUpAllExecutableOnEnv("PATH") },
	func() []Element { return listUpAllExecutableOnEnv("NYAGOSPATH") },
}

// AppendCommandLister is the function to append the environment variable name at seeing on command-name completion.
func AppendCommandLister(f func() []Element) {
	commandListUpper = append(commandListUpper, f)
}

// HookToList is the slice for Completion-Hook functions for users.
var HookToList = []func(context.Context, *readline.Buffer, *List) (*List, error){}

func init() {
	f := readline.KeyGoFuncT{Func: KeyFuncCompletion, Name: "COMPLETE"}
	err := readline.BindKeyFunc(readline.K_CTRL_I, &f)
	if err != nil {
		panic(err.Error())
	}
}
