package lua

import "fmt"
import "os"
import "strings"

import "../alias/table"
import "../interpreter"
import "../option"

func alias(L *Lua) int {
	name := L.ToString(1)
	value := L.ToString(2)
	aliasTable.Table[strings.ToLower(name)] = value
	return 0
}

func setEnv(L *Lua) int {
	name := L.ToString(1)
	value := L.ToString(2)
	os.Setenv(name, value)
	return 0
}

func getEnv(L *Lua) int {
	name := L.ToString(1)
	value := os.Getenv(name)
	if len(value) > 0 {
		L.PushString(value)
		return 1
	} else {
		return 0
	}
}

func exec(L *Lua) int {
	statement := L.ToString(1)
	_, err := interpreter.Interpret(statement, option.CommandHooks, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return 0
}

func SetFunctions(this *Lua) {
	stackPos := this.GetTop()
	this.NewTable()
	this.PushGoFunction(alias)
	this.SetField(-2, "alias")
	this.PushGoFunction(setEnv)
	this.SetField(-2, "setenv")
	this.PushGoFunction(getEnv)
	this.SetField(-2, "getenv")
	this.PushGoFunction(exec)
	this.SetField(-2, "exec")
	this.SetGlobal("nyagos")

	// replace io.getenv
	this.GetGlobal("io")
	this.PushGoFunction(getEnv)
	this.SetField(-2, "getenv")

	this.SetTop(stackPos)
}
