// +build !vanilla

package mains

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/yuin/gopher-lua"
	"github.com/zetamatta/nyagos/shell"
)

type argsFilterFlagT struct{}

var argsFilterFlag argsFilterFlagT

var orgArgHook func(context.Context, *shell.Shell, []string) ([]string, error)

func newArgHook(ctx context.Context, it *shell.Shell, args []string) ([]string, error) {
	if ctx.Value(argsFilterFlag) != nil {
		return orgArgHook(ctx, it, args)
	}
	ctx = context.WithValue(ctx, argsFilterFlag, true)

	luawrapper, ok := it.Tag().(*luaWrapper)
	if !ok {
		return nil, errors.New("Could not get lua instance(newArgHook)")
	}
	L := luawrapper.Lua
	nyagosTable := L.GetGlobal("nyagos")
	if _, ok := nyagosTable.(*lua.LTable); !ok {
		return orgArgHook(ctx, it, args)
	}
	f := L.GetField(nyagosTable, "argsfilter")
	if _, ok := f.(*lua.LFunction); !ok {
		return orgArgHook(ctx, it, args)
	}
	param := L.NewTable()
	for i := 0; i < len(args); i++ {
		L.SetTable(param, lua.LNumber(i), lua.LString(args[i]))
	}
	L.Push(f)
	L.Push(param)
	if err := callLua(ctx, it, 1, 1); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return orgArgHook(ctx, it, args)
	}
	resultTmp := L.Get(-1)
	L.Pop(1)
	result, ok := resultTmp.(*lua.LTable)
	if !ok {
		return orgArgHook(ctx, it, args)
	}
	size := result.Len()
	newargs := make([]string, size+1)
	for i := 0; i <= size; i++ {
		newargs[i] = L.GetTable(result, lua.LNumber(i)).String()
	}
	return orgArgHook(ctx, it, newargs)
}
