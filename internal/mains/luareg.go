package mains

import (
	"context"

	"github.com/yuin/gopher-lua"
)

const (
	readlineLuaRegistryKey = "nyagos.readline"
	shellLuaRegistryKey    = "nyagos.shell"
)

func getLuaRegistry(L Lua, key string) any {
	tbl, ok := L.Get(lua.RegistryIndex).(*lua.LTable)
	if !ok {
		return nil
	}
	ud, ok := L.GetField(tbl, key).(*lua.LUserData)
	if !ok {
		return nil
	}
	return ud.Value
}

func setLuaRegistry(L Lua, key string, value any) {
	reg, ok := L.Get(lua.RegistryIndex).(*lua.LTable)
	if !ok {
		return
	}
	if value != nil {
		ud := L.NewUserData()
		ud.Value = value
		L.SetField(reg, key, ud)
	} else {
		L.SetField(reg, key, lua.LNil)
	}
}

func pushLuaRegistry(L Lua, key string, value any) func() {
	orig := getLuaRegistry(L, key)
	setLuaRegistry(L, key, value)
	return func() {
		setLuaRegistry(L, key, orig)
	}
}

func setContext(ctx context.Context, L Lua) {
	if ctx == nil {
		ctx = context.Background()
	}
	L.SetContext(ctx)
}

func getContext(L Lua) context.Context {
	return L.Context()
}
