package main

import (
	"github.com/yuin/gopher-lua"
	"github.com/zetamatta/nyagos/functions"
)

type Lua = *lua.LState

func NewLua() (Lua, error) {
	L := lua.NewState()

	nyagosTable := L.NewTable()

	for name, function := range functions.Table {
		L.SetTable(nyagosTable, lua.LString(name), L.NewFunction(lua2cmd(function)))
	}
	L.SetGlobal("nyagos", nyagosTable)

	shareTable := L.NewTable()
	L.SetGlobal("share", shareTable)

	return L, nil
}

func luaArgsToInterfaces(L Lua) []interface{} {
	end := L.GetTop()
	var param []interface{}
	if end > 0 {
		param = make([]interface{}, 0, end-1)
		for i := 1; i <= end; i++ {
			switch L.Get(i).Type() {
			case lua.LTString:
				param = append(param, L.ToString(i))
			case lua.LTNumber:
				param = append(param, L.ToInt(i))
			default:
				param = append(param, lua.LNil)
			}
		}
	} else {
		param = []interface{}{}
	}
	return param
}

func pushInterfaces(L Lua, values []interface{}) {
	for _, valueTmp := range values {
		if valueTmp == nil {
			L.Push(lua.LNil)
		} else {
			switch value := valueTmp.(type) {
			case string:
				L.Push(lua.LString(value))
			case int:
				L.Push(lua.LNumber(value))
			}
		}
	}
}

func lua2cmd(f func([]interface{}) []interface{}) func(Lua) int {
	return func(L Lua) int {
		param := luaArgsToInterfaces(L)
		result := f(param)
		pushInterfaces(L, result)
		return len(result)
	}
}
