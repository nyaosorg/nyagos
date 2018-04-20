package main

import (
	"github.com/yuin/gopher-lua"
)

func deepCopyTable(L2 *lua.LState, t1 *lua.LTable, t2 *lua.LTable) {
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

func cloneTo(L1, L2 *lua.LState) bool {
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
	return true
}

func Clone(L *lua.LState) *lua.LState {
	L2 := lua.NewState()
	if cloneTo(L, L2) {
		return L2
	} else {
		L2.Close()
		return nil
	}
}
