package gopherSh

import (
	"context"
	"fmt"
	"os"

	"github.com/yuin/gopher-lua"
	"github.com/zetamatta/nyagos/shell"
)

type LuaFilterStream struct {
	shell.Stream
	L Lua
}

func (this *LuaFilterStream) ReadLine(ctx context.Context) (context.Context, string, error) {
	ctx, line, err := this.Stream.ReadLine(ctx)
	if err != nil {
		return ctx, "", err
	}

	L := this.L

	stackPos := L.GetTop()
	defer L.SetTop(stackPos)

	nyagosTable, ok := L.GetGlobal("nyagos").(*lua.LTable)
	if !ok {
		return ctx, line, err
	}
	luaFilter, ok := L.GetField(nyagosTable, "filter").(*lua.LFunction)
	if !ok {
		return ctx, line, err
	}

	L.Push(luaFilter)
	L.Push(lua.LString(line))
	defer setContext(L, getContext(L))
	setContext(L, ctx)
	err = L.PCall(1, 1, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ctx, line, nil
	}
	newLine, ok := L.Get(-1).(lua.LString)
	L.Pop(1)
	if !ok {
		return ctx, line, nil
	}
	return ctx, string(newLine), nil
}
