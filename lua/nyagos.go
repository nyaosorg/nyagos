package lua

import "os"
import "strings"

import "../alias/table"

func LuaAlias(L *Lua) int {
	name := L.ToString(1)
	value := L.ToString(2)
	aliasTable.Table[strings.ToLower(name)] = value
	return 0
}

func LuaSetEnv(L *Lua) int {
	name := L.ToString(1)
	value := L.ToString(2)
	os.Setenv(name, value)
	return 0
}

func SetFunctions(this *Lua) {
	this.PushGoFunction(LuaAlias)
	this.SetGlobal("alias")
	this.PushGoFunction(LuaSetEnv)
	this.SetGlobal("setenv")
}
