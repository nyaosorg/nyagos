//go:build !vanilla
// +build !vanilla

package mains

import (
	"context"
	"fmt"
	"os"

	"github.com/nyaosorg/nyagos/internal/shell"
	"github.com/yuin/gopher-lua"
)

type luaFilterStream struct {
	shell.Stream
	L Lua
}

func luaLineFilter(ctx context.Context, L Lua, line string) string {
	stackPos := L.GetTop()
	defer L.SetTop(stackPos)

	nyagosTable, ok := L.GetGlobal("nyagos").(*lua.LTable)
	if !ok {
		return line
	}
	luaFilter, ok := L.GetField(nyagosTable, "filter").(*lua.LFunction)
	if !ok {
		return line
	}

	L.Push(luaFilter)
	L.Push(lua.LString(line))
	defer setContext(getContext(L), L)
	setContext(ctx, L)
	err := L.PCall(1, 1, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return line
	}
	newLine, ok := L.Get(-1).(lua.LString)
	L.Pop(1)
	if !ok {
		return line
	}
	return string(newLine)
}

func (lfs *luaFilterStream) ReadLine(ctx context.Context) (string, error) {
	line, err := lfs.Stream.ReadLine(ctx)
	if err != nil {
		return "", err
	}
	return luaLineFilter(ctx, lfs.L, line), err
}
