package mains

import (
	"context"
	"fmt"
	"os"

	"github.com/zetamatta/nyagos/mains/lua-dll"
	"github.com/zetamatta/nyagos/shell"
)

var luaFilter lua.Object = lua.TNil{}

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

	L.Push(luaFilter)
	if !L.IsFunction(-1) {
		return ctx, line, nil
	}
	L.PushString(line)
	err = L.CallWithContext(ctx, 1, 1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ctx, line, nil
	}
	if !L.IsString(-1) {
		return ctx, line, nil
	}
	newLine, err := L.ToString(-1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ctx, line, nil
	}
	return ctx, newLine, nil
}
