package main

import (
	"context"
	"errors"
	"strings"

	"github.com/yuin/gopher-lua"
	"github.com/zetamatta/nyagos/alias"
	"github.com/zetamatta/nyagos/shell"
)

type LuaBinaryChank struct {
	Chank *lua.LFunction
}

func (this *LuaBinaryChank) String() string {
	return "(lua-function)"
}

func (this *LuaBinaryChank) Call(ctx context.Context, cmd *shell.Cmd) (int, error) {
	luawrapper, ok := cmd.Tag().(*luaWrapper)
	if !ok {
		return 255, errors.New("LuaBinaryChank.Call: Lua instance not found")
	}
	L := luawrapper.Lua
	ctx = context.WithValue(ctx, luaKey, L)
	L.Push(this.Chank)

	table := L.NewTable()
	for i, arg1 := range cmd.Args() {
		L.SetTable(table, lua.LNumber(i), lua.LString(arg1))
	}
	L.Push(table)

	callLua(ctx, &cmd.Shell, 1, 0)

	return 1, nil
}

func cmdSetAlias(L Lua) int {
	key := strings.ToLower(L.ToString(-2))
	switch L.Get(-1).Type() {
	case lua.LTString:
		alias.Table[key] = alias.New(L.ToString(-1))
	case lua.LTFunction:
		alias.Table[key] = &LuaBinaryChank{Chank: L.ToFunction(-1)}
	case lua.LTNil:
		delete(alias.Table, key)
	}
	L.Push(lua.LTrue)
	return 1
}

func cmdGetAlias(L Lua) int {
	value, ok := alias.Table[strings.ToLower(L.ToString(-1))]
	if !ok {
		L.Push(lua.LNil)
		return 1
	}
	switch v := value.(type) {
	case *LuaBinaryChank:
		L.Push(v.Chank)
	default:
		L.Push(lua.LString(v.String()))
	}
	return 1
}
