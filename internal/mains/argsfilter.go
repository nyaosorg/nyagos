//go:build !vanilla
// +build !vanilla

package mains

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nyaosorg/nyagos/internal/shell"
	"github.com/yuin/gopher-lua"
)

var orgArgHook func(context.Context, *shell.Shell, []string, []string) ([]string, []string, error)

func newArgHook(ctx context.Context, it *shell.Shell, args, rawargs []string) ([]string, []string, error) {
	saveHook := it.ArgsHook
	it.ArgsHook = orgArgHook
	defer func() {
		it.ArgsHook = saveHook
	}()

	luawrapper, ok := it.Tag().(*luaWrapper)
	if !ok {
		return nil, nil, errors.New("Could not get lua instance(newArgHook)")
	}
	L := luawrapper.Lua
	nyagosTable := L.GetGlobal("nyagos")
	if _, ok := nyagosTable.(*lua.LTable); !ok {
		return orgArgHook(ctx, it, args, rawargs)
	}
	f := L.GetField(nyagosTable, "argsfilter")
	if _, ok := f.(*lua.LFunction); !ok {
		return orgArgHook(ctx, it, args, rawargs)
	}
	param := L.NewTable()
	for i := 0; i < len(args); i++ {
		L.SetTable(param, lua.LNumber(i), lua.LString(args[i]))
	}
	L.Push(f)
	L.Push(param)
	if err := execLuaKeepContextAndShell(ctx, it, L, 1, 1); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return orgArgHook(ctx, it, args, rawargs)
	}
	resultTmp := L.Get(-1)
	L.Pop(1)
	result, ok := resultTmp.(*lua.LTable)
	if !ok {
		return orgArgHook(ctx, it, args, rawargs)
	}
	size := result.Len()
	newargs := make([]string, size+1)
	newrawargs := make([]string, size+1)
	for i := 0; i <= size; i++ {
		arg1 := L.GetTable(result, lua.LNumber(i)).String()
		newargs[i] = arg1
		if strings.ContainsAny(arg1, ` "&|<>`) {
			arg1 = shell.Quote(arg1)
		}
		newrawargs[i] = arg1
	}
	return orgArgHook(ctx, it, newargs, newrawargs)
}
