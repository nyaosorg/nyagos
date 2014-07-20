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

func exec(L *Lua) int {
	statement := L.ToString(1)
	_, err := interpreter.Interpret(statement, option.CommandHooks, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return 0
}

func SetFunctions(this *Lua) {
	this.PushGoFunction(alias)
	this.SetGlobal("alias")
	this.PushGoFunction(setEnv)
	this.SetGlobal("setenv")
	this.PushGoFunction(exec)
	this.SetGlobal("exec")
}
