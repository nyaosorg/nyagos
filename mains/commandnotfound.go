package main

import (
	"errors"

	"github.com/zetamatta/nyagos/mains/lua-dll"
	"github.com/zetamatta/nyagos/shell"
)

var orgOnCommandNotFound func(*shell.Cmd, error) error

var luaOnCommandNotFound lua.Object = lua.TNil{}

func onCommandNotFound(sh *shell.Cmd, err error) error {
	luawrapper, ok := sh.Tag().(*luaWrapper)
	if !ok {
		return errors.New("Could get lua instance(on_command_not_found)")
	}
	L := luawrapper.Lua

	L.Push(luaOnCommandNotFound)
	if !L.IsFunction(-1) {
		L.Pop(1)
		return orgOnCommandNotFound(sh, err)
	}
	L.NewTable()
	for key, val := range sh.Args() {
		L.PushString(val)
		L.RawSetI(-2, lua.Integer(key))
	}
	err1 := L.Call(1, 1)
	defer L.Pop(1)
	if err1 != nil {
		return err
	}
	if L.ToBool(-1) {
		return nil
	} else {
		return orgOnCommandNotFound(sh, err)
	}
}
