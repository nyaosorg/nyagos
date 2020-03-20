package completion

import (
	"context"
	"github.com/zetamatta/go-readline-ny"
)

var commandListUpper = []func(context.Context) ([]Element, error){
	func(ctx context.Context) ([]Element, error) {
		return listUpAllExecutableOnEnv(ctx, "PATH")
	},
	func(ctx context.Context) ([]Element, error) {
		return listUpAllExecutableOnEnv(ctx, "NYAGOSPATH")
	},
}

// AppendCommandLister is the function to append the environment variable name at seeing on command-name completion.
func AppendCommandLister(f func(context.Context) ([]Element, error)) {
	commandListUpper = append(commandListUpper, f)
}

// HookToList is the slice for Completion-Hook functions for users.
var HookToList = []func(context.Context, *readline.Buffer, *List) (*List, error){}

func init() {
	f := readline.KeyGoFuncT{Func: KeyFuncCompletion, Name: "COMPLETE"}
	err := readline.GlobalKeyMap.BindKeyFunc(readline.K_CTRL_I, &f)
	if err != nil {
		panic(err.Error())
	}
}
