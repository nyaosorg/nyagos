package mains

import (
	"github.com/yuin/gopher-lua"
)

func bit32and(L *lua.LState) int {
	result := ^0
	for n := L.GetTop(); n > 0; n-- {
		value, ok := L.Get(n).(lua.LNumber)
		if !ok {
			L.Push(lua.LNil)
			L.Push(lua.LString("bit32.and: found NaN"))
			return 2
		}
		result = result & int(value)
	}
	L.Push(lua.LNumber(result))
	return 1
}

func SetupBit32Table(L *lua.LState) {
	table := L.NewTable()
	L.SetField(table, "band", L.NewFunction(bit32and))
	L.SetGlobal("bit32", table)
}
