package nua

import "fmt"
import "os"
import "strings"
import "os/exec"

import . "../lua"
import "../alias"
import "../interpreter"

type LuaFunction struct {
	L            *Lua
	registoryKey string
}

func (this LuaFunction) String() string {
	return "<<Lua-function>>"
}

func (this LuaFunction) Call(cmd *exec.Cmd) (interpreter.NextT, error) {
	this.L.GetField(Registory, this.registoryKey)
	for _, arg1 := range cmd.Args {
		this.L.PushString(arg1)
	}
	err := this.L.Call(len(cmd.Args), 0)
	return interpreter.CONTINUE, err
}

func cmdAlias(L *Lua) int {
	name := L.ToString(1)
	key := strings.ToLower(name)
	if L.IsString(2) {
		value := L.ToString(2)
		alias.Table[key] = alias.New(value)
	} else if L.IsFunction(2) {
		regkey := "nyagos.alias." + key
		L.SetField(Registory, regkey)
		alias.Table[key] = LuaFunction{L, regkey}
	}
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

func cmdEcho(L *Lua) int {
	n := L.GetTop()
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Print("\t")
		}
		fmt.Print(L.ToString(i))
	}
	fmt.Print("\n")
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
	this.PushGoFunction(cmdEcho)
	this.SetField(-2, "echo")
	this.SetGlobal("nyagos")

	// replace io.getenv
	this.GetGlobal("os")
	this.PushGoFunction(cmdGetEnv)
	this.SetField(-2, "getenv")

	this.SetTop(stackPos)
}
