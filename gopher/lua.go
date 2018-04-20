package main

import (
	"github.com/yuin/gopher-lua"
)

type Lua = *lua.LState

func NewLua() (Lua, error) {
	this := lua.NewState()

	nyagosTable := this.NewTable()
	this.SetGlobal("nyagos", nyagosTable)

	shareTable := this.NewTable()
	this.SetGlobal("share", shareTable)

	return this, nil
}
