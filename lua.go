package main

import "fmt"
import "io"
import "os"
import "os/exec"
import "strings"
import "unsafe"

import . "./lua"
import "./alias"
import "./interpreter"
import "./mbcs"

const nyagos_exec_cmd = "nyagos.exec.cmd"

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
	this.L.PushLightUserData(unsafe.Pointer(cmd))
	this.L.SetField(Registory, nyagos_exec_cmd)
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
	L.GetField(Registory, nyagos_exec_cmd)
	cmd := (*exec.Cmd)(L.ToUserData(-1))
	L.Pop(1)
	var out io.Writer
	if cmd != nil {
		out = cmd.Stdout
	} else {
		out = os.Stdout
	}

	n := L.GetTop()
	for i := 1; i <= n; i++ {
		if i > 1 {
			fmt.Fprint(out, "\t")
		}
		fmt.Fprint(out, L.ToString(i))
	}
	fmt.Fprint(out, "\n")
	return 0
}

func cmdAtoU(L *Lua) int {
	L.PushString(mbcs.AtoU(L.ToAnsiString(1)))
	return 1
}

func cmdUtoA(L *Lua) int {
	L.PushAnsiString(mbcs.UtoA(L.ToString(1)))
	return 1
}

func SetLuaFunctions(this *Lua) {
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
	this.PushGoFunction(cmdAtoU)
	this.SetField(-2, "atou")
	this.PushGoFunction(cmdUtoA)
	this.SetField(-2, "utoa")
	this.SetGlobal("nyagos")

	// replace io.getenv
	this.GetGlobal("os")
	this.PushGoFunction(cmdGetEnv)
	this.SetField(-2, "getenv")

	this.SetTop(stackPos)
}
