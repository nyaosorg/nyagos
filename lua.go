package main

import "fmt"
import "os"

import "./dos"
import "./interpreter"
import "./lua"

type LuaFunction struct {
	L            *lua.Lua
	registoryKey string
}

var LuaInstanceToCmd = map[uintptr]*interpreter.Interpreter{}

func (this LuaFunction) String() string {
	return "<<Lua-function>>"
}

func (this LuaFunction) Call(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	this.L.GetField(lua.REGISTORYINDEX, this.registoryKey)
	this.L.NewTable()
	for i, arg1 := range cmd.Args {
		this.L.PushInteger(i)
		this.L.PushString(arg1)
		this.L.SetTable(-3)
	}
	LuaInstanceToCmd[this.L.Id()] = cmd
	err := this.L.Call(1, 0)
	return interpreter.CONTINUE, err
}

func SetLuaFunctions(this *lua.Lua) {
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
	this.PushGoFunction(cmdWrite)
	this.SetField(-2, "write")
	this.PushGoFunction(cmdAtoU)
	this.SetField(-2, "atou")
	this.PushGoFunction(cmdUtoA)
	this.SetField(-2, "utoa")
	this.PushGoFunction(cmdGetwd)
	this.SetField(-2, "getwd")
	this.PushGoFunction(cmdWhich)
	this.SetField(-2, "which")
	this.PushGoFunction(cmdEval)
	this.SetField(-2, "eval")
	this.PushGoFunction(cmdGlob)
	this.SetField(-2, "glob")
	this.PushGoFunction(cmdBindKey)
	this.SetField(-2, "bindkey")
	this.PushGoFunction(cmdGetHistory)
	this.SetField(-2, "gethistory")
	exeName, exeNameErr := dos.GetModuleFileName()
	if exeNameErr != nil {
		fmt.Fprintln(os.Stderr, exeNameErr)
	} else {
		this.PushString(exeName)
		this.SetField(-2, "exe")
	}
	this.SetGlobal("nyagos")

	// replace io.getenv
	this.GetGlobal("os")
	this.PushGoFunction(cmdGetEnv)
	this.SetField(-2, "getenv")

	var orgArgHook func([]string) []string
	orgArgHook = interpreter.SetArgsHook(func(args []string) []string {
		pos := this.GetTop()
		defer this.SetTop(pos)
		this.GetGlobal("nyagos")
		this.GetField(-1, "argsfilter")
		if !this.IsFunction(-1) {
			return orgArgHook(args)
		}
		this.NewTable()
		for i := 0; i < len(args); i++ {
			this.PushInteger(i)
			this.PushString(args[i])
			this.SetTable(-3)
		}
		if err := this.Call(1, 1); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return orgArgHook(args)
		}
		if this.GetType(-1) != lua.TTABLE {
			return orgArgHook(args)
		}
		newargs := []string{}
		for i := 0; true; i++ {
			this.PushInteger(i)
			this.GetTable(-2)
			if this.GetType(-1) == lua.TNIL {
				break
			}
			arg1, arg1err := this.ToString(-1)
			if arg1err == nil {
				newargs = append(newargs, arg1)
			} else {
				fmt.Fprintln(os.Stderr, arg1err.Error())
			}
			this.Pop(1)
		}
		return orgArgHook(newargs)
	})
}
