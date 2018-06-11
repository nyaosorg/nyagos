package mains

import (
	"github.com/yuin/gopher-lua"
	"github.com/zetamatta/nyagos/functions"
)

func SetupBit32Table(L *lua.LState) {
	table := L.NewTable()
	L.SetField(table, "band", L.NewFunction(lua2cmd(functions.CmdBitAnd)))
	L.SetField(table, "bor", L.NewFunction(lua2cmd(functions.CmdBitOr)))
	L.SetField(table, "bxor", L.NewFunction(lua2cmd(functions.CmdBitXor)))
	L.SetGlobal("bit32", table)
}
