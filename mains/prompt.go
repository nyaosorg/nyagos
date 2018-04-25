package mains

import (
	"context"
	"fmt"
	"os"

	"github.com/zetamatta/nyagos/functions"
	"github.com/zetamatta/nyagos/mains/lua-dll"
	"github.com/zetamatta/nyagos/shell"
)

var prompt_hook lua.Object = lua.TGoFunction(lua2cmd(functions.Prompt))

func printPrompt(ctx context.Context, sh *shell.Shell, L Lua) (int, error) {
	L.Push(prompt_hook)

	if !L.IsFunction(-1) {
		L.Pop(1)
		return 0, nil
	}
	L.PushString(os.Getenv("PROMPT"))
	if err := callCSL(ctx, sh, L, 1, 1); err != nil {
		return 0, err
	}
	length, lengthErr := L.ToInteger(-1)
	L.Pop(1)
	if lengthErr == nil {
		return length, nil
	} else {
		return 0, fmt.Errorf("nyagos.prompt: return-value(length) is invalid: %s", lengthErr.Error())
	}
}
