package main

import "bytes"
import "fmt"
import "io"
import "os"
import "os/exec"
import "strings"
import "unsafe"

import "./alias"
import "./conio/readline"
import "./dos"
import "./interpreter"
import "./lua"
import "./mbcs"

import "github.com/shiena/ansicolor"

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

func cmdAlias(L *lua.Lua) int {
	name, nameErr := L.ToString(1)
	if nameErr != nil {
		L.PushNil()
		L.PushString(nameErr.Error())
		return 2
	}
	key := strings.ToLower(name)
	switch L.GetType(2) {
	case lua.TSTRING:
		value, err := L.ToString(2)
		if err == nil {
			alias.Table[key] = alias.New(value)
		} else {
			L.PushNil()
			L.PushString(err.Error())
			return 2
		}
	case lua.TFUNCTION:
		regkey := "nyagos.alias." + key
		L.SetField(lua.REGISTORYINDEX, regkey)
		alias.Table[key] = LuaFunction{L, regkey}
	}
	L.PushBool(true)
	return 1
}

func cmdSetEnv(L *lua.Lua) int {
	name, nameErr := L.ToString(1)
	if nameErr != nil {
		L.PushNil()
		L.PushString(nameErr.Error())
		return 2
	}
	value, valueErr := L.ToString(2)
	if valueErr != nil {
		L.PushNil()
		L.PushString(valueErr.Error())
		return 2
	}
	os.Setenv(name, value)
	L.PushBool(true)
	return 1
}

func cmdGetEnv(L *lua.Lua) int {
	name, nameErr := L.ToString(1)
	if nameErr != nil {
		L.PushNil()
		return 1
	}
	value := os.Getenv(name)
	if len(value) > 0 {
		L.PushString(value)
	} else {
		L.PushNil()
	}
	return 1
}

func cmdExec(L *lua.Lua) int {
	statement, statementErr := L.ToString(1)
	if statementErr != nil {
		L.PushNil()
		L.PushString(statementErr.Error())
		return 2
	}
	_, err := interpreter.New().Interpret(statement)

	if err != nil {
		L.PushNil()
		L.PushString(err.Error())
		return 2
	}
	L.PushBool(true)
	return 1
}

func cmdEval(L *lua.Lua) int {
	statement, statementErr := L.ToString(1)
	if statementErr != nil {
		L.PushNil()
		L.PushString(statementErr.Error())
		return 2
	}
	r, w, err := os.Pipe()
	if err != nil {
		L.PushNil()
		L.PushString(err.Error())
		return 2
	}
	go func(statement string, w *os.File) {
		it := interpreter.New()
		it.Stdout = w
		it.Interpret(statement)
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
	L.PushAnsiString(result)
	return 1
}

func cmdWrite(L *lua.Lua) int {
	var out io.Writer = os.Stdout
	cmd, cmdOk := LuaInstanceToCmd[L.Id()]
	if cmdOk && cmd != nil && cmd.Stdout != nil {
		out = cmd.Stdout
	}
	switch out.(type) {
	case *os.File:
		out = ansicolor.NewAnsiColorWriter(out)
	}

	n := L.GetTop()
	for i := 1; i <= n; i++ {
		str, err := L.ToString(i)
		if err != nil {
			L.PushNil()
			L.PushString(err.Error())
			return 2
		}
		if i > 1 {
			fmt.Fprint(out, "\t")
		}
		fmt.Fprint(out, str)
	}
	L.PushBool(true)
	return 1
}

func cmdGetwd(L *lua.Lua) int {
	wd, err := os.Getwd()
	if err == nil {
		L.PushString(wd)
		return 1
	} else {
		return 0
	}
}

func cmdWhich(L *lua.Lua) int {
	if L.GetType(-1) != lua.TSTRING {
		return 0
	}
	name, nameErr := L.ToString(-1)
	if nameErr != nil {
		L.PushNil()
		L.PushString(nameErr.Error())
		return 2
	}
	path, err := exec.LookPath(name)
	if err == nil {
		L.PushString(path)
		return 1
	} else {
		L.PushNil()
		L.PushString(err.Error())
		return 2
	}
}

func cmdAtoU(L *lua.Lua) int {
	str, err := mbcs.AtoU(L.ToAnsiString(1))
	if err == nil {
		L.PushString(str)
		return 1
	} else {
		return 0
	}
}

func cmdUtoA(L *lua.Lua) int {
	utf8, utf8err := L.ToString(1)
	if utf8err != nil {
		L.PushNil()
		L.PushString(utf8err.Error())
		return 2
	}
	str, err := mbcs.UtoA(utf8)
	if err == nil {
		if len(str) >= 1 {
			L.PushAnsiString(str[:len(str)-1])
		} else {
			L.PushString("")
		}
		L.PushNil()
		return 2
	} else {
		L.PushNil()
		L.PushString(err.Error())
		return 2
	}
}

func cmdGlob(L *lua.Lua) int {
	if !L.IsString(-1) {
		return 0
	}
	wildcard, wildcardErr := L.ToString(-1)
	if wildcardErr != nil {
		L.PushNil()
		L.PushString(wildcardErr.Error())
		return 2
	}
	list, err := dos.Glob(wildcard)
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

type KeyLuaFuncT struct {
	L            *lua.Lua
	registoryKey string
}

func getBufferForCallBack(L *lua.Lua) (*readline.Buffer, int) {
	if L.GetType(1) != lua.TTABLE {
		L.PushNil()
		L.PushString("bindKeyExec: call with : not .")
		return nil, 2
	}
	L.GetField(1, "buffer")
	if L.GetType(-1) != lua.TLIGHTUSERDATA {
		L.PushNil()
		L.PushString("bindKey.Call: invalid object")
		return nil, 2
	}
	buffer := (*readline.Buffer)(L.ToUserData(-1))
	if buffer == nil {
		L.PushNil()
		L.PushString("bindKey.Call: invalid member")
		return nil, 2
	}
	return buffer, 0
}

func callInsert(L *lua.Lua) int {
	buffer, stackRc := getBufferForCallBack(L)
	if buffer == nil {
		return stackRc
	}
	text, textErr := L.ToString(2)
	if textErr != nil {
		L.PushNil()
		L.PushString(textErr.Error())
		return 2
	}
	buffer.InsertAndRepaint(text)
	L.PushBool(true)
	return 1
}

func callKeyFunc(L *lua.Lua) int {
	buffer, stackRc := getBufferForCallBack(L)
	if buffer == nil {
		return stackRc
	}
	key, keyErr := L.ToString(2)
	if keyErr != nil {
		L.PushNil()
		L.PushString(keyErr.Error())
		return 2
	}
	function, funcErr := readline.GetFunc(key)
	if funcErr != nil {
		L.PushNil()
		L.PushString(funcErr.Error())
		return 2
	}
	rc := function.Call(buffer)
	L.PushBool(true)
	switch rc {
	case readline.ENTER:
		L.PushBool(true)
		return 2
	case readline.ABORT:
		L.PushBool(false)
		return 2
	}
	return 1
}

func (this *KeyLuaFuncT) Call(buffer *readline.Buffer) readline.Result {
	this.L.GetField(lua.REGISTORYINDEX, this.registoryKey)
	this.L.NewTable()
	pos := -1
	var text bytes.Buffer
	for i, c := range buffer.Buffer {
		if i == buffer.Cursor {
			pos = text.Len() + 1
		}
		text.WriteRune(c)
	}
	if pos < 0 {
		pos = text.Len()
	}
	this.L.PushInteger(pos)
	this.L.SetField(-2, "pos")
	this.L.PushString(text.String())
	this.L.SetField(-2, "text")
	this.L.PushLightUserData(unsafe.Pointer(buffer))
	this.L.SetField(-2, "buffer")
	this.L.PushGoFunction(callKeyFunc)
	this.L.SetField(-2, "call")
	this.L.PushGoFunction(callInsert)
	this.L.SetField(-2, "insert")
	if err := this.L.Call(1, 1); err != nil {
		fmt.Println(os.Stderr, err)
	}
	switch this.L.GetType(-1) {
	case lua.TSTRING:
		str, strErr := this.L.ToString(-1)
		if strErr == nil {
			buffer.InsertAndRepaint(str)
		}
	case lua.TBOOLEAN:
		if !this.L.ToBool(-1) {
			buffer.Buffer = []rune{}
			buffer.Length = 0
		}
		return readline.ENTER
	}
	return readline.CONTINUE
}

func cmdBindKey(L *lua.Lua) int {
	key, keyErr := L.ToString(-2)
	if keyErr != nil {
		L.PushString(keyErr.Error())
		return 1
	}
	key = strings.Replace(strings.ToUpper(key), "-", "_", -1)
	switch L.GetType(-1) {
	case lua.TFUNCTION:
		regkey := "nyagos.bind." + key
		L.SetField(lua.REGISTORYINDEX, regkey)
		if err := readline.BindKeyFunc(key, &KeyLuaFuncT{L, regkey}); err != nil {
			L.PushNil()
			L.PushString(err.Error())
			return 2
		} else {
			L.PushBool(true)
			return 1
		}
	default:
		val, valErr := L.ToString(-1)
		if valErr != nil {
			L.PushNil()
			L.PushString(valErr.Error())
			return 2
		}
		err := readline.BindKeySymbol(key, val)
		if err != nil {
			L.PushNil()
			L.PushString(err.Error())
			return 2
		} else {
			L.PushBool(true)
			return 1
		}
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
