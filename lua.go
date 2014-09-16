package main

import "fmt"
import "io"
import "os"
import "os/exec"
import "strings"
import "unsafe"

import . "./lua"
import "./alias"
import "./dos"
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
	_, err := interpreter.Interpret(statement, nil)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	return 0
}

func cmdEval(L *Lua) int {
	statement := L.ToString(1)
	r, w, err := os.Pipe()
	if err != nil {
		return 0
	}
	go func(statement string, w *os.File) {
		interpreter.Interpret(statement, &interpreter.Stdio{Stdout: w})
		w.Close()
	}(statement, w)

	var result = []byte{}
	for {
		buffer := make([]byte, 256)
		size, err := r.Read(buffer)
		if err != nil || size <= 0 {
			break
		}
		result = append(result, buffer[0:size]...)
	}
	r.Close()
	if result != nil {
		L.PushAnsiString(result)
		return 1
	} else {
		return 0
	}
}

func cmdEcho(L *Lua) int {
	var out io.Writer
	L.GetField(Registory, nyagos_exec_cmd)
	if L.GetType(-1) == TLIGHTUSERDATA {
		cmd := (*exec.Cmd)(L.ToUserData(-1))
		L.Pop(1)
		if cmd != nil {
			out = cmd.Stdout
		} else {
			out = os.Stdout
		}
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

func cmdWhich(L *Lua) int {
	if L.GetType(-1) != TSTRING {
		return 0
	}
	name := L.ToString(-1)
	path, err := exec.LookPath(name)
	if err == nil {
		L.PushString(path)
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

func cmdGlob(L *Lua) int {
	if !L.IsString(-1) {
		return 0
	}
	list, err := dos.Glob(L.ToString(-1))
	if err != nil {
		L.PushNil()
		L.PushString(err.Error())
		return 2
	} else {
		L.NewTable()
		for i := 0; i < len(list); i++ {
			L.PushInteger(i + 1)
			L.PushString(list[i])
			L.SetTable(-3)
		}
		return 1
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
	this.PushGoFunction(cmdWhich)
	this.SetField(-2, "which")
	this.PushGoFunction(cmdEval)
	this.SetField(-2, "eval")
	this.PushGoFunction(cmdGlob)
	this.SetField(-2, "glob")
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
		if this.GetType(-1) != TTABLE {
			return orgArgHook(args)
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
		return orgArgHook(newargs)
	})
}
