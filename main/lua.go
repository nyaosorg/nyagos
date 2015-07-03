package main

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"../dos"
	"../interpreter"
	"../lua"
)

type LuaFunction struct {
	L            lua.Lua
	registoryKey string
}

var LuaInstanceToCmd = map[uintptr]*interpreter.Interpreter{}

func (this LuaFunction) String() string {
	return "<<Lua-function>>"
}

var mutex4dll sync.Mutex
var luaUsedOnThatPipeline = map[uint]uint{}

const ERRMSG_CAN_NOT_USE_TWO_LUA_ON = "Can not use two Lua-command on the same pipeline"
const ERRMSG_CAN_NOT_RUN_LUA_ON_BACKGROUND = "Can not run Lua-Command on background"

func (this LuaFunction) Call(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	if cmd.IsBackGround {
		fmt.Fprintf(os.Stderr, "%s: %s\n",
			cmd.Args[0],
			ERRMSG_CAN_NOT_RUN_LUA_ON_BACKGROUND)
		return interpreter.CONTINUE,
			errors.New(ERRMSG_CAN_NOT_RUN_LUA_ON_BACKGROUND)
	}
	L := this.L
	seq := cmd.PipeSeq
	mutex4dll.Lock()
	if p, ok := luaUsedOnThatPipeline[seq[0]]; ok && p != seq[1] {
		mutex4dll.Unlock()
		fmt.Fprintf(os.Stderr, "%s: %s\n",
			cmd.Args[0],
			ERRMSG_CAN_NOT_USE_TWO_LUA_ON)
		return interpreter.CONTINUE, errors.New(ERRMSG_CAN_NOT_USE_TWO_LUA_ON)
	}
	luaUsedOnThatPipeline[seq[0]] = seq[1]
	mutex4dll.Unlock()

	L.GetField(lua.LUA_REGISTRYINDEX, this.registoryKey)
	L.NewTable()
	for i, arg1 := range cmd.Args {
		L.PushString(arg1)
		L.RawSetI(-2, lua.Integer(i))
	}
	save := LuaInstanceToCmd[L.State()]
	LuaInstanceToCmd[L.State()] = cmd
	err := L.Call(1, 1)
	if err == nil {
		newargs := make([]string, 0)
		if L.IsTable(-1) {
			L.PushInteger(0)
			L.GetTable(-2)
			if val, err1 := L.ToString(-1); val != "" && err1 == nil {
				newargs = append(newargs, val)
			}
			L.Pop(1)
			for i := 1; ; i++ {
				L.PushInteger(lua.Integer(i))
				L.GetTable(-2)
				if L.IsNil(-1) {
					L.Pop(1)
					break
				}
				val, err1 := L.ToString(-1)
				L.Pop(1)
				if err1 != nil {
					break
				}
				newargs = append(newargs, val)
			}
			it := cmd.Clone()
			it.Args = newargs
			it.Spawnvp()
		} else if val, err1 := L.ToString(-1); val != "" && err1 == nil {
			cmd.Clone().Interpret(val)
		}
	}
	L.Pop(1)
	LuaInstanceToCmd[this.L.State()] = save
	return interpreter.CONTINUE, err
}

const original_io_lines = "original_io_lines"

func ioLines(this lua.Lua) int {
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

func ioLinesNext(this lua.Lua) int {
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

var orgArgHook func(*interpreter.Interpreter, []string) []string

func newArgHook(it *interpreter.Interpreter, args []string) []string {
	L, Lok := it.Tag.(lua.Lua)
	if !Lok {
		panic("main/lua.go: can get interpreter instance")
	}
	pos := L.GetTop()
	defer L.SetTop(pos)
	L.GetGlobal("nyagos")
	L.GetField(-1, "argsfilter")
	if !L.IsFunction(-1) {
		return orgArgHook(it, args)
	}
	L.NewTable()
	for i := 0; i < len(args); i++ {
		L.PushString(args[i])
		L.RawSetI(-2, lua.Integer(i))
	}
	if err := L.Call(1, 1); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return orgArgHook(it, args)
	}
	if L.GetType(-1) != lua.LUA_TTABLE {
		return orgArgHook(it, args)
	}
	newargs := []string{}
	for i := lua.Integer(0); true; i++ {
		L.PushInteger(i)
		L.GetTable(-2)
		if L.GetType(-1) == lua.LUA_TNIL {
			break
		}
		arg1, arg1err := L.ToString(-1)
		if arg1err == nil {
			newargs = append(newargs, arg1)
		} else {
			fmt.Fprintln(os.Stderr, arg1err.Error())
		}
		L.Pop(1)
	}
	return orgArgHook(it, newargs)
}

var orgOnCommandNotFound func(*interpreter.Interpreter, error) error

func on_command_not_found(inte *interpreter.Interpreter, err error) error {
	L, Lok := inte.Tag.(lua.Lua)
	if !Lok {
		panic("on_command_not_found: Interpreter.Tag is not lua instance")
	}
	L.GetGlobal("nyagos")
	L.GetField(-1, "on_command_not_found")
	L.Remove(-2) // remove nyagos.
	if L.IsFunction(-1) {
		L.NewTable()
		for key, val := range inte.Args {
			L.PushString(val)
			L.RawSetI(-2, lua.Integer(key))
		}
		L.Call(1, 1)
		defer L.Pop(1)
		if L.ToBool(-1) {
			return nil
		} else {
			return orgOnCommandNotFound(inte, err)
		}
	} else {
		L.Pop(1)
		return orgOnCommandNotFound(inte, err)
	}
}

func SetLuaFunctions(this lua.Lua) {
	stackPos := this.GetTop()
	defer this.SetTop(stackPos)

	nyagos_table := map[string]interface{}{
		"access":       cmdAccess,
		"atou":         cmdAtoU,
		"bindkey":      cmdBindKey,
		"commonprefix": cmdCommonPrefix,
		"eval":         cmdEval,
		"exec":         cmdExec,
		"getalias":     cmdGetAlias,
		"getenv":       cmdGetEnv,
		"gethistory":   cmdGetHistory,
		"getkey":       cmdGetKey,
		"getviewwidth": cmdGetViewWidth,
		"getwd":        cmdGetwd,
		"glob":         cmdGlob,
		"pathjoin":     cmdPathJoin,
		"raweval":      cmdRawEval,
		"setalias":     cmdSetAlias,
		"setenv":       cmdSetEnv,
		"setrunewidth": cmdSetRuneWidth,
		"shellexecute": cmdShellExecute,
		"stat":         cmdStat,
		"utoa":         cmdUtoA,
		"which":        cmdWhich,
		"write":        cmdWrite,
		"writerr":      cmdWriteErr,
	}
	if exeName, exeNameErr := dos.GetModuleFileName(); exeNameErr != nil {
		fmt.Fprintln(os.Stderr, exeNameErr)
	} else {
		nyagos_table["exe"] = exeName
	}
	this.Push(nyagos_table)

	this.NewTable() // "nyagos.alias"
	this.NewTable() // metatable.
	this.PushGoFunction(cmdSetAlias)
	this.SetField(-2, "__call")
	this.PushGoFunction(cmdSetAlias)
	this.SetField(-2, "__newindex")
	this.PushGoFunction(cmdGetAlias)
	this.SetField(-2, "__index")
	this.SetMetaTable(-2)
	this.SetField(-2, "alias")

	this.NewTable() // "nyagos.env"
	this.NewTable() // metatable
	this.PushGoFunction(cmdSetEnv)
	this.SetField(-2, "__newindex")
	this.PushGoFunction(cmdGetEnv)
	this.SetField(-2, "__index")
	this.SetMetaTable(-2)
	this.SetField(-2, "env")

	this.SetGlobal("nyagos")

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

	orgArgHook = interpreter.SetArgsHook(newArgHook)

	orgOnCommandNotFound = interpreter.OnCommandNotFound
	interpreter.OnCommandNotFound = on_command_not_found
}
