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
		this.L.PushString(arg1)
		this.L.RawSetI(-2, i)
	}
	LuaInstanceToCmd[this.L.State()] = cmd
	err := this.L.Call(1, 0)
	return interpreter.CONTINUE, err
}

const original_io_lines = "original_io_lines"

func ioLines(this *lua.Lua) int {
	if this.IsString(1) {
		// io.lines("FILENAME") --> use original io.lines
		this.GetField(lua.REGISTORYINDEX, original_io_lines)
		this.PushValue(1)
		this.Call(1, 1)
	} else {
		// io.lines() --> use nyagos version
		this.PushGoFunction(ioLinesNext)
	}
	return 1
}

func ioLinesNext(this *lua.Lua) int {
	cmd := LuaInstanceToCmd[this.State()]

	line := make([]byte, 0, 256)
	var ch [1]byte
	for {
		n, err := cmd.Stdin.Read(ch[0:1])
		if n <= 0 || err != nil {
			if len(line) <= 0 {
				this.PushNil()
			} else {
				this.PushAnsiString(line)
			}
			return 1
		}
		if ch[0] == '\n' {
			this.PushAnsiString(line)
			return 1
		}
		line = append(line, ch[0])
	}
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
	this.PushGoFunction(cmdAccess)
	this.SetField(-2, "access")
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
	this.PushGoFunction(cmdSetRuneWidth)
	this.SetField(-2, "setrunewidth")
	this.PushGoFunction(cmdShellExecute)
	this.SetField(-2, "shellexecute")
	exeName, exeNameErr := dos.GetModuleFileName()
	if exeNameErr != nil {
		fmt.Fprintln(os.Stderr, exeNameErr)
	} else {
		this.PushString(exeName)
		this.SetField(-2, "exe")
	}
	this.SetGlobal("nyagos")

	// replace os.getenv
	this.GetGlobal("os")           // +1
	this.PushGoFunction(cmdGetEnv) // +2
	this.SetField(-2, "getenv")    // +1
	this.Pop(1)                    // 0

	// save io.lines as original_io_lines
	this.GetGlobal("io")                                 // +1
	this.GetField(-1, "lines")                           // +2
	this.SetField(lua.REGISTORYINDEX, original_io_lines) // +1
	this.Pop(1)                                          // 0

	// replace io.lines
	this.GetGlobal("io")         // +1
	this.PushGoFunction(ioLines) // +2
	this.SetField(-2, "lines")   // +1
	this.Pop(1)                  // 0

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
			this.PushString(args[i])
			this.RawSetI(-2, i)
		}
		if err := this.Call(1, 1); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return orgArgHook(args)
		}
		if this.GetType(-1) != lua.LUA_TTABLE {
			return orgArgHook(args)
		}
		newargs := []string{}
		for i := 0; true; i++ {
			this.PushInteger(i)
			this.GetTable(-2)
			if this.GetType(-1) == lua.LUA_TNIL {
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
