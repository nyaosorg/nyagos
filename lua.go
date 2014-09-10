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
	this.L.NewTable()
	for i, arg1 := range cmd.Args {
		this.L.PushInteger(i)
		this.L.PushString(arg1)
		this.L.SetTable(-3)
	}
	this.L.PushLightUserData(unsafe.Pointer(cmd))
	this.L.SetField(Registory, nyagos_exec_cmd)
	err := this.L.Call(1, 0)
	return interpreter.CONTINUE, err
}

func cmdAlias(L *Lua) int {
	name := L.ToString(1)
	key := strings.ToLower(name)
	switch L.GetType(2) {
	case TSTRING:
		value := L.ToString(2)
		alias.Table[key] = alias.New(value)
	case TFUNCTION:
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
	if L.GetType(-1) != TUSERDATE {
		fmt.Fprintln(os.Stderr, "nyagos.echo: invalid argument: not userdata")
		return 0
	}
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

func cmdGetwd(L *Lua) int {
	wd, err := os.Getwd()
	if err == nil {
		L.PushString(wd)
		return 1
	} else {
		return 0
	}
}

func cmdAtoU(L *Lua) int {
	str, err := mbcs.AtoU(L.ToAnsiString(1))
	if err == nil {
		L.PushString(str)
		return 1
	} else {
		return 0
	}
}

func cmdUtoA(L *Lua) int {
	str, err := mbcs.UtoA(L.ToString(1))
	if err == nil {
		L.PushAnsiString(str)
		return 1
	} else {
		return 0
	}
}

func SetLuaFunctions(this *Lua) {
	stackPos := this.GetTop()
	defer this.SetTop(stackPos)
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
	this.PushGoFunction(cmdGetwd)
	this.SetField(-2, "getwd")
	this.SetGlobal("nyagos")

	// replace io.getenv
	this.GetGlobal("os")
	this.PushGoFunction(cmdGetEnv)
	this.SetField(-2, "getenv")

	interpreter.ArgsHook = func(args []string) []string {
		pos := this.GetTop()
		defer this.SetTop(pos)
		this.GetGlobal("nyagos")
		this.GetField(-1, "argsfilter")
		if !this.IsFunction(-1) {
			return args
		}
		this.NewTable()
		for i := 0; i < len(args); i++ {
			this.PushInteger(i)
			this.PushString(args[i])
			this.SetTable(-3)
		}
		if err := this.Call(1, 1); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return args
		}
		if this.GetType(-1) != TTABLE {
			return args
		}
		newargs := []string{}
		for i := 0; true; i++ {
			this.PushInteger(i)
			this.GetTable(-2)
			if this.GetType(-1) == TNIL {
				break
			}
			newargs = append(newargs, this.ToString(-1))
			this.Pop(1)
		}
		return newargs
	}
}
