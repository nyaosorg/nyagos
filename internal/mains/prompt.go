//go:build !vanilla
// +build !vanilla

package mains

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/nyaosorg/nyagos/internal/functions"
	"github.com/nyaosorg/nyagos/internal/shell"
	"github.com/yuin/gopher-lua"
)

func printPrompt(ctx context.Context, sh *shell.Shell, L Lua, w io.Writer) (int, error) {
	nyagosTbl := L.GetGlobal("nyagos")
	prompt := L.GetField(nyagosTbl, "prompt")
	if promptHook, ok := prompt.(*lua.LFunction); ok {
		// nyagos.prompt is function.
		L.Push(promptHook)
		L.Push(lua.LString(os.Getenv("PROMPT")))
		if err := execLuaKeepContextAndShell(ctx, sh, L, 1, 1); err != nil {
			return 0, err
		}
		defer L.Pop(1)
		if promptString, ok := L.Get(-1).(lua.LString); ok {
			return io.WriteString(w, string(promptString))
		} else if length, ok := L.Get(-1).(lua.LNumber); ok {
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
	return io.WriteString(w, functions.PromptCore(w, promptStr))
}
