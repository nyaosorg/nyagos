package mains

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/zetamatta/nyagos/mains/lua-dll"
	"github.com/zetamatta/nyagos/shell"
)

var orgArgHook func(context.Context, *shell.Shell, []string) ([]string, error)

var luaArgsFilter lua.Object = lua.TNil{}

func newArgHook(ctx context.Context, it *shell.Shell, args []string) ([]string, error) {
	luawrapper, ok := it.Tag().(*luaWrapper)
	if !ok {
		return nil, errors.New("Could not get lua instance(newArgHook)")
	}
	L := luawrapper.Lua
	L.Push(luaArgsFilter)
	if !L.IsFunction(-1) {
		L.Pop(1)
		return orgArgHook(ctx, it, args)
	}
	L.NewTable()
	for i := 0; i < len(args); i++ {
		L.PushString(args[i])
		L.RawSetI(-2, lua.Integer(i))
	}
	if err := callLua(ctx, it, 1, 1); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return orgArgHook(ctx, it, args)
	}
	if L.GetType(-1) != lua.LUA_TTABLE {
		return orgArgHook(ctx, it, args)
	}
	newargs := []string{}
	for i := lua.Integer(0); true; i++ {
		L.PushInteger(i)
		L.GetTable(-2)
		if L.GetType(-1) == lua.LUA_TNIL {
			break
		}
		arg1, arg1err := L.ToString(-1)
		if arg1err == nil {
			newargs = append(newargs, arg1)
		} else {
			fmt.Fprintln(os.Stderr, arg1err.Error())
		}
		L.Pop(1)
	}
	return orgArgHook(ctx, it, newargs)
}
