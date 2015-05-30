package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/shiena/ansicolor"

	"../alias"
	"../completion"
	"../conio"
	"../dos"
	"../interpreter"
	"../lua"
)

const alias_prefix = "nyagos.alias."

func cmdSetAlias(L *lua.Lua) int {
	name, nameErr := L.ToString(1)
	if nameErr != nil {
		return L.Push(nil, nameErr)
	}
	key := strings.ToLower(name)
	switch L.GetType(2) {
	case lua.LUA_TSTRING:
		value, err := L.ToString(2)
		regkey := alias_prefix + key
		L.SetField(lua.LUA_REGISTRYINDEX, regkey)
		if err == nil {
			alias.Table[key] = alias.New(value)
		} else {
			return L.Push(nil, err)
		}
	case lua.LUA_TFUNCTION:
		regkey := alias_prefix + key
		L.SetField(lua.LUA_REGISTRYINDEX, regkey)
		alias.Table[key] = LuaFunction{L, regkey}
	}
	return L.Push(true)
}

func cmdGetAlias(L *lua.Lua) int {
	name, nameErr := L.ToString(1)
	if nameErr != nil {
		return L.Push(nil, nameErr)
	}
	regkey := alias_prefix + strings.ToLower(name)
	L.GetField(lua.LUA_REGISTRYINDEX, regkey)
	return 1
}

func cmdSetEnv(L *lua.Lua) int {
	name, nameErr := L.ToString(1)
	if nameErr != nil {
		return L.Push(nil, nameErr)
	}
	value, valueErr := L.ToString(2)
	if valueErr != nil {
		return L.Push(nil, valueErr)
	}
	if len(value) > 0 {
		os.Setenv(name, value)
	} else {
		os.Unsetenv(name)
	}
	return L.Push(true)
}

func cmdGetEnv(L *lua.Lua) int {
	name, nameErr := L.ToString(1)
	if nameErr != nil {
		return L.Push(nil)
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
	var err error
	if L.IsTable(1) {
		L.Len(1)
		n, _ := L.ToInteger(-1)
		L.Pop(1)
		args := make([]string, 0, n+1)
		for i := 0; i <= n; i++ {
			L.PushInteger(lua.Integer(i))
			L.GetTable(-2)
			arg1, err := L.ToString(-1)
			if err == nil && arg1 != "" {
				args = append(args, arg1)
			}
			L.Pop(1)
		}
		interpreter1 := interpreter.New()
		interpreter1.Args = args
		_, err = interpreter1.Spawnvp()
	} else {
		statement, statementErr := L.ToString(1)
		if statementErr != nil {
			return L.Push(nil, statementErr)
		}
		_, err = interpreter.New().Interpret(statement)
	}
	if err != nil {
		var out io.Writer = os.Stderr
		if cmd, cmdOk := LuaInstanceToCmd[L.State()]; cmdOk {
			out = cmd.Stderr
		}
		fmt.Fprintln(out, err)
		return L.Push(nil, err)
	}
	return L.Push(true)
}

type emptyWriter struct{}

func (e *emptyWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func cmdEval(L *lua.Lua) int {
	statement, statementErr := L.ToString(1)
	if statementErr != nil {
		return L.Push(nil, statementErr)
	}
	r, w, err := os.Pipe()
	if err != nil {
		return L.Push(nil, err)
	}
	go func(statement string, w *os.File) {
		it := interpreter.New()
		it.SetStdout(w)
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
	L.PushAnsiString(bytes.Trim(result, "\r\n\t "))
	return 1
}

func cmdWrite(L *lua.Lua) int {
	var out io.Writer = os.Stdout
	cmd, cmdOk := LuaInstanceToCmd[L.State()]
	if cmdOk && cmd != nil && cmd.Stdout != nil {
		out = cmd.Stdout
	}
	return cmdWriteSub(L, out)
}

func cmdWriteErr(L *lua.Lua) int {
	var out io.Writer = os.Stderr
	cmd, cmdOk := LuaInstanceToCmd[L.State()]
	if cmdOk && cmd != nil && cmd.Stderr != nil {
		out = cmd.Stderr
	}
	return cmdWriteSub(L, out)
}

func cmdWriteSub(L *lua.Lua, out io.Writer) int {
	switch out.(type) {
	case *os.File:
		out = ansicolor.NewAnsiColorWriter(out)
	}
	n := L.GetTop()
	for i := 1; i <= n; i++ {
		str, err := L.ToString(i)
		if err != nil {
			return L.Push(nil, err)
		}
		if i > 1 {
			fmt.Fprint(out, "\t")
		}
		fmt.Fprint(out, str)
	}
	return L.Push(true)
}

func cmdGetwd(L *lua.Lua) int {
	wd, err := dos.Getwd()
	if err == nil {
		return L.Push(wd)
	} else {
		return L.Push(nil, err)
	}
}

func cmdWhich(L *lua.Lua) int {
	if L.GetType(-1) != lua.LUA_TSTRING {
		return 0
	}
	name, nameErr := L.ToString(-1)
	if nameErr != nil {
		return L.Push(nil, nameErr)
	}
	path, err := exec.LookPath(name)
	if err == nil {
		return L.Push(path)
	} else {
		return L.Push(nil, err)
	}
}

func cmdAtoU(L *lua.Lua) int {
	str, err := dos.AtoU(L.ToAnsiString(1))
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
		return L.Push(nil, utf8err)
	}
	str, err := dos.UtoA(utf8)
	if err != nil {
		return L.Push(nil, err)
	}
	if len(str) >= 1 {
		L.PushAnsiString(str[:len(str)-1])
	} else {
		L.PushString("")
	}
	L.PushNil()
	return 2
}

func cmdGlob(L *lua.Lua) int {
	result := make([]string, 0)
	for i := 1; ; i++ {
		wildcard, wildcardErr := L.ToString(i)
		if wildcard == "" || wildcardErr != nil {
			break
		}
		list, err := dos.Glob(wildcard)
		if list == nil || err != nil {
			result = append(result, wildcard)
		} else {
			result = append(result, list...)
		}
	}
	L.NewTable()
	for i := 0; i < len(result); i++ {
		L.PushString(result[i])
		L.RawSetI(-2, lua.Integer(i+1))
	}
	return 1
}

func cmdGetHistory(this *lua.Lua) int {
	if this.GetType(-1) == lua.LUA_TNUMBER {
		val, err := this.ToInteger(-1)
		if err != nil {
			return this.Push(nil, err.Error())
		}
		this.PushString(conio.DefaultEditor.GetHistoryAt(val).Line)
	} else {
		this.PushInteger(lua.Integer(conio.DefaultEditor.HistoryLen()))
	}
	return 1
}

func cmdSetRuneWidth(this *lua.Lua) int {
	char, charErr := this.ToInteger(1)
	if charErr != nil {
		return this.Push(nil, charErr)
	}
	width, widthErr := this.ToInteger(2)
	if widthErr != nil {
		return this.Push(nil, widthErr)
	}
	conio.SetCharWidth(rune(char), width)
	this.PushBool(true)
	return 1
}

func cmdShellExecute(this *lua.Lua) int {
	action, actionErr := this.ToString(1)
	if actionErr != nil {
		return this.Push(nil, actionErr)
	}
	path, pathErr := this.ToString(2)
	if pathErr != nil {
		return this.Push(nil, pathErr)
	}
	param, paramErr := this.ToString(3)
	if paramErr != nil {
		param = ""
	}
	dir, dirErr := this.ToString(4)
	if dirErr != nil {
		dir = ""
	}
	err := dos.ShellExecute(action, path, param, dir)
	if err != nil {
		return this.Push(nil, err)
	} else {
		return this.Push(true)
	}
}

func cmdAccess(L *lua.Lua) int {
	path, pathErr := L.ToString(1)
	if pathErr != nil {
		return L.Push(nil, pathErr)
	}
	mode, modeErr := L.ToInteger(2)
	if modeErr != nil {
		return L.Push(nil, modeErr)
	}
	fi, err := os.Stat(path)

	var result bool
	if err != nil || fi == nil {
		result = false
	} else {
		switch {
		case mode == 0:
			result = true
		case mode&1 != 0: // X_OK
		case mode&2 != 0: // W_OK
			result = fi.Mode().Perm()&0200 != 0
		case mode&4 != 0: // R_OK
			result = fi.Mode().Perm()&0400 != 0
		}
	}
	L.PushBool(result)
	return 1
}

func cmdPathJoin(L *lua.Lua) int {
	path, pathErr := L.ToString(1)
	if pathErr != nil {
		return L.Push(nil, pathErr)
	}
	for i, i_ := 2, L.GetTop(); i <= i_; i++ {
		pathI, pathIErr := L.ToString(i)
		if pathIErr != nil {
			return L.Push(nil, pathErr)
		}
		path = dos.Join(path, pathI)
	}
	return L.Push(path, nil)
}

func cmdCommonPrefix(L *lua.Lua) int {
	if L.GetType(1) != lua.LUA_TTABLE {
		return 0
	}
	list := []string{}
	for i := lua.Integer(1); true; i++ {
		L.PushInteger(i)
		L.GetTable(1)
		if str, err := L.ToString(2); err == nil && str != "" {
			list = append(list, str)
		} else {
			break
		}
		L.Remove(2)
	}
	L.PushString(completion.CommonPrefix(list))
	return 1
}

func cmdGetKey(L *lua.Lua) int {
	keycode, scancode, shiftstatus := conio.GetKey()
	L.PushInteger(lua.Integer(keycode))
	L.PushInteger(lua.Integer(scancode))
	L.PushInteger(lua.Integer(shiftstatus))
	return 3
}

func cmdGetViewWidth(L *lua.Lua) int {
	width, _ := conio.GetScreenBufferInfo().ViewSize()
	L.PushInteger(lua.Integer(width))
	return 1
}
