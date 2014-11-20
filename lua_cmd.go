package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/shiena/ansicolor"

	"./alias"
	"./conio"
	"./dos"
	"./history"
	"./interpreter"
	"./lua"
)

func cmdAlias(L *lua.Lua) int {
	name, nameErr := L.ToString(1)
	if nameErr != nil {
		L.PushNil()
		L.PushString(nameErr.Error())
		return 2
	}
	key := strings.ToLower(name)
	switch L.GetType(2) {
	case lua.LUA_TSTRING:
		value, err := L.ToString(2)
		if err == nil {
			alias.Table[key] = alias.New(value)
		} else {
			L.PushNil()
			L.PushString(err.Error())
			return 2
		}
	case lua.LUA_TFUNCTION:
		regkey := "nyagos.alias." + key
		L.SetField(lua.LUA_REGISTRYINDEX, regkey)
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

type emptyWriter struct{}

func (e *emptyWriter) Write(b []byte) (int, error) {
	return len(b), nil
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
		it.Stderr = &emptyWriter{}
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
	if L.GetType(-1) != lua.LUA_TSTRING {
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
		L.PushNil()
		L.PushString(utf8err.Error())
		return 2
	}
	str, err := dos.UtoA(utf8)
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
			L.PushString(list[i])
			L.RawSetI(-2, i+1)
		}
		return 1
	}
}

func cmdGetHistory(this *lua.Lua) int {
	if this.GetType(-1) == lua.LUA_TNUMBER {
		val, err := this.ToInteger(-1)
		if err != nil {
			this.PushNil()
			this.PushString(err.Error())
		}
		this.PushString(history.Get(val))
	} else {
		this.PushInteger(history.Len())
	}
	return 1
}

func cmdSetRuneWidth(this *lua.Lua) int {
	char, charErr := this.ToInteger(1)
	if charErr != nil {
		this.PushNil()
		this.PushString(charErr.Error())
		return 2
	}
	width, widthErr := this.ToInteger(2)
	if widthErr != nil {
		this.PushNil()
		this.PushString(widthErr.Error())
		return 2
	}
	conio.SetCharWidth(rune(char), width)
	this.PushBool(true)
	return 1
}

func cmdShellExecute(this *lua.Lua) int {
	action, actionErr := this.ToString(1)
	if actionErr != nil {
		this.PushNil()
		this.PushString(actionErr.Error())
		return 2
	}
	path, pathErr := this.ToString(2)
	if pathErr != nil {
		this.PushNil()
		this.PushString(pathErr.Error())
		return 2
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
		this.PushNil()
		this.PushString(err.Error())
	} else {
		this.PushBool(true)
	}
	return 1
}

func cmdAccess(L *lua.Lua) int {
	path, pathErr := L.ToString(1)
	if pathErr != nil {
		L.PushNil()
		L.PushString(pathErr.Error())
		return 2
	}
	mode, modeErr := L.ToInteger(2)
	if modeErr != nil {
		L.PushNil()
		L.PushString(modeErr.Error())
		return 2
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
