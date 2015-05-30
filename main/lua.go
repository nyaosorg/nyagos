package main

import (
	"fmt"
	"os"

	"../dos"
	"../interpreter"
	"../lua"
)

type LuaFunction struct {
	L            *lua.Lua
	registoryKey string
}

var LuaInstanceToCmd = map[uintptr]*interpreter.Interpreter{}

func (this LuaFunction) String() string {
	return "<<Lua-function>>"
}

func (this LuaFunction) Call(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	this.L.GetField(lua.LUA_REGISTRYINDEX, this.registoryKey)
	this.L.NewTable()
	for i, arg1 := range cmd.Args {
		this.L.PushString(arg1)
		this.L.RawSetI(-2, lua.Integer(i))
	}
	LuaInstanceToCmd[this.L.State()] = cmd
	err := this.L.Call(1, 0)
	return interpreter.CONTINUE, err
}

const original_io_lines = "original_io_lines"

func ioLines(this *lua.Lua) int {
	if this.IsString(1) {
		// io.lines("FILENAME") --> use original io.lines
		this.GetField(lua.LUA_REGISTRYINDEX, original_io_lines)
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

	nyagos_table := map[string]interface{}{
		"access":       cmdAccess,
		"setalias":     cmdSetAlias,
		"atou":         cmdAtoU,
		"commonprefix": cmdCommonPrefix,
		"bindkey":      cmdBindKey,
		"eval":         cmdEval,
		"exec":         cmdExec,
		"getalias":     cmdGetAlias,
		"getenv":       cmdGetEnv,
		"gethistory":   cmdGetHistory,
		"getkey":       cmdGetKey,
		"getwd":        cmdGetwd,
		"glob":         cmdGlob,
		"pathjoin":     cmdPathJoin,
		"setenv":       cmdSetEnv,
		"setrunewidth": cmdSetRuneWidth,
		"shellexecute": cmdShellExecute,
		"utoa":         cmdUtoA,
		"which":        cmdWhich,
		"write":        cmdWrite,
		"writerr":      cmdWriteErr,
		"getviewwidth": cmdGetViewWidth,
	}
	if exeName, exeNameErr := dos.GetModuleFileName(); exeNameErr != nil {
		fmt.Fprintln(os.Stderr, exeNameErr)
	} else {
		nyagos_table["exe"] = exeName
	}
	this.Push(nyagos_table)
	this.SetGlobal("nyagos")

	if err := this.LoadString(`
		nyagos.alias = setmetatable({},{
			__call=function(t,k,v) nyagos.setalias(k,v) end,
			__newindex=function(t,k,v) nyagos.setalias(k,v) end,
			__index=function(t,k) return nyagos.getalias(k,v) end
		})
		nyagos.env = setmetatable({},{
			__newindex=function(t,k,v) nyagos.setenv(k,v) end,
			__index=function(t,k) return nyagos.getenv(k) end
		})
	`); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else if err = this.Call(0, 0); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	// replace os.getenv
	this.GetGlobal("os")           // +1
	this.PushGoFunction(cmdGetEnv) // +2
	this.SetField(-2, "getenv")    // +1
	this.Pop(1)                    // 0

	// save io.lines as original_io_lines
	this.GetGlobal("io")                                    // +1
	this.GetField(-1, "lines")                              // +2
	this.SetField(lua.LUA_REGISTRYINDEX, original_io_lines) // +1
	this.Pop(1)                                             // 0

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
			this.RawSetI(-2, lua.Integer(i))
		}
		if err := this.Call(1, 1); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return orgArgHook(args)
		}
		if this.GetType(-1) != lua.LUA_TTABLE {
			return orgArgHook(args)
		}
		newargs := []string{}
		for i := lua.Integer(0); true; i++ {
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
