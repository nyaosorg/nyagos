//go:build !vanilla
// +build !vanilla

package mains

import (
	"errors"

	"github.com/yuin/gopher-lua"
)

func deepCopyTable(L2 Lua, t1 *lua.LTable, t2 *lua.LTable) {
	t1.ForEach(func(key, val lua.LValue) {
		target := L2.GetTable(t2, key)
		if target.Type() == lua.LTNil { // do not override
			if tbl1, ok := val.(*lua.LTable); ok {
				tbl2 := L2.NewTable()
				deepCopyTable(L2, tbl1, tbl2)
				L2.SetTable(t2, key, tbl2)
			} else {
				L2.SetTable(t2, key, val)
			}
		}
	})
}

func cloneTo(L1, L2 Lua) bool {
	G1, ok := L1.GetGlobal("_G").(*lua.LTable)
	if !ok {
		return false
	}
	G2, ok := L2.GetGlobal("_G").(*lua.LTable)
	if !ok {
		L2.Close()
		return false
	}
	deepCopyTable(L2, G1, G2)

	reg1, ok := L1.Get(lua.RegistryIndex).(*lua.LTable)
	if !ok {
		return false
	}
	reg2, ok := L2.Get(lua.RegistryIndex).(*lua.LTable)
	if !ok {
		L2.Close()
		return false
	}
	deepCopyTable(L2, reg1, reg2)

	if ctx := L1.Context(); ctx != nil {
		L2.SetContext(ctx)
	}
	return true
}

// Clone makes a copy of Lua instance.
func Clone(L Lua) (Lua, error) {
	L2, err := NewLua()
	if err != nil {
		return L2, err
	}
	if cloneTo(L, L2) {
		return L2, nil
	}
	L2.Close()
	return nil, errors.New("could not create Lua instance")
}
