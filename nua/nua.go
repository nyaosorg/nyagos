package nua

import "fmt"
import "os"
import "strings"

import . "../lua"
import "../alias"
import "../interpreter"

func cmdAlias(L *Lua) int {
	name := L.ToString(1)
	value := L.ToString(2)
	alias.Table[strings.ToLower(name)] = value
	return 0
}

func cmdSetEnv(L *Lua) int {
	name := L.ToString(1)
	value := L.ToString(2)
	os.Setenv(name, value)
	return 0
}

func cmdGetEnv(L *Lua) int {
	name := L.ToString(1)
	value := os.Getenv(name)
	if len(value) > 0 {
		L.PushString(value)
		return 1
	} else {
		return 0
	}
}

func cmdExec(L *Lua) int {
	statement := L.ToString(1)
	_, err := interpreter.Interpret(statement, alias.Hook, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return 0
}

func SetFunctions(this *Lua) {
	stackPos := this.GetTop()
	this.NewTable()
	this.PushGoFunction(cmdAlias)
	this.SetField(-2, "alias")
	this.PushGoFunction(cmdSetEnv)
	this.SetField(-2, "setenv")
	this.PushGoFunction(cmdGetEnv)
	this.SetField(-2, "getenv")
	this.PushGoFunction(cmdExec)
	this.SetField(-2, "exec")
	this.SetGlobal("nyagos")

	// replace io.getenv
	this.GetGlobal("os")
	this.PushGoFunction(cmdGetEnv)
	this.SetField(-2, "getenv")

	this.SetTop(stackPos)
}
