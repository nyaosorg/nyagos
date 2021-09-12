//go:build !vanilla
// +build !vanilla

package mains

import (
	"context"
	"errors"
	"os"

	"github.com/nyaosorg/nyagos/functions"
	"github.com/nyaosorg/nyagos/shell"
	"github.com/yuin/gopher-lua"
)

func printPrompt(ctx context.Context, sh *shell.Shell, L Lua) (int, error) {
	nyagosTbl := L.GetGlobal("nyagos")
	prompt := L.GetField(nyagosTbl, "prompt")
	if promptHook, ok := prompt.(*lua.LFunction); ok {
		// nyagos.prompt is function.
		L.Push(promptHook)
		L.Push(lua.LString(os.Getenv("PROMPT")))
		if err := execLuaKeepContextAndShell(ctx, sh, L, 1, 1); err != nil {
			return 0, err
		}

		length, ok := L.Get(-1).(lua.LNumber)
		L.Pop(1)
		if ok {
			return int(length), nil
		}
		return 0, errors.New("nyagos.prompt: return-value(length) is not a number")
	}
	var promptStr string
	if promptLStr, ok := prompt.(lua.LString); ok {
		promptStr = string(promptLStr)
	} else {
		promptStr = os.Getenv("PROMPT")
	}
	return functions.PromptCore(sh.Term(), promptStr), nil
}
