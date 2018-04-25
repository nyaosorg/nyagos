package main

import (
	"context"
	"errors"
	"os"

	"github.com/yuin/gopher-lua"
	"github.com/zetamatta/nyagos/shell"
)

func printPrompt(ctx context.Context, sh *shell.Shell, L Lua) (int, error) {
	nyagosTbl := L.GetGlobal("nyagos")
	promptHook, ok := L.GetField(nyagosTbl, "prompt").(*lua.LFunction)

	if !ok {
		return 0, nil
	}

	L.Push(promptHook)
	L.Push(lua.LString(os.Getenv("PROMPT")))
	if err := callCSL(ctx, sh, L, 1, 1); err != nil {
		return 0, err
	}

	length, ok := L.Get(-1).(lua.LNumber)
	L.Pop(1)
	if ok {
		return int(length), nil
	} else {
		return 0, errors.New("nyagos.prompt: return-value(length) is not a number")
	}
}
