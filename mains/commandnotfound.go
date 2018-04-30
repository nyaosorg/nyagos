package main

import (
	"context"
	"errors"

	"github.com/zetamatta/nyagos/mains/lua-dll"
	"github.com/zetamatta/nyagos/shell"
)

var orgOnCommandNotFound func(context.Context, *shell.Cmd, error) error

var luaOnCommandNotFound lua.Object = lua.TNil{}

func onCommandNotFound(ctx context.Context, sh *shell.Cmd, err error) error {
	luawrapper, ok := sh.Tag().(*luaWrapper)
	if !ok {
		return errors.New("Could get lua instance(on_command_not_found)")
	}
	L := luawrapper.Lua

	L.Push(luaOnCommandNotFound)
	if !L.IsFunction(-1) {
		L.Pop(1)
		return orgOnCommandNotFound(ctx, sh, err)
	}
	L.NewTable()
	for key, val := range sh.Args() {
		L.PushString(val)
		L.RawSetI(-2, lua.Integer(key))
	}
	err1 := callLua(ctx, &sh.Shell, 1, 1)
	defer L.Pop(1)
	if err1 != nil {
		return err
	}
	if L.ToBool(-1) {
		return nil
	} else {
		return orgOnCommandNotFound(ctx, sh, err)
	}
}
