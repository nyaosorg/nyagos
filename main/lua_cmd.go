package main

import (
	"bytes"
	"errors"
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
	"../dos/ansicfile"
	"../interpreter"
	"../lua"
)

type LuaBinaryChank struct {
	Chank []byte
}

func (this *LuaBinaryChank) String() string {
	return "(lua-function)"
}

func (this *LuaBinaryChank) Call(cmd *interpreter.Interpreter) (interpreter.ErrorLevel, error) {
	L, L_ok := cmd.Tag.(lua.Lua)
	if !L_ok {
		return interpreter.ErrorLevel(255), errors.New("LuaBinaryChank.Call: Lua instance not found")
	}
	if err := L.LoadBufferX(cmd.Args[0], this.Chank, "b"); err != nil {
		return interpreter.ErrorLevel(255), err
	}
	L.NewTable()
	for i, arg1 := range cmd.Args {
		L.PushString(arg1)
		L.RawSetI(-2, lua.Integer(i))
	}
	L.NewTable()
	for i, arg1 := range cmd.RawArgs {
		L.PushString(arg1)
		L.RawSetI(-2, lua.Integer(i))
	}
	L.SetField(-2, "rawargs")
	err := NyagosCallLua(L, cmd, 1, 1)
	errorlevel := interpreter.NOERROR
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
			it, err1 := cmd.Clone()
			if err1 != nil {
				errorlevel = interpreter.ErrorLevel(255)
				err = err1
			} else {
				it.Args = newargs
				errorlevel, err = it.Spawnvp()
			}
		} else if val, err1 := L.ToInteger(-1); err1 == nil {
			errorlevel = interpreter.ErrorLevel(val)
		} else if val, err1 := L.ToString(-1); val != "" && err1 == nil {
			it, err1 := cmd.Clone()
			if err1 != nil {
				errorlevel = interpreter.ErrorLevel(255)
				err = err1
			} else {
				errorlevel, err = it.Interpret(val)
			}
		}
	}
	L.Pop(1)
	return errorlevel, err
}

func cmdSetAlias(L lua.Lua) int {
	name, nameErr := L.ToString(-2)
	if nameErr != nil {
		return L.Push(nil, nameErr.Error())
	}
	key := strings.ToLower(name)
	switch L.GetType(-1) {
	case lua.LUA_TSTRING:
		value, err := L.ToString(-1)
		if err == nil {
			alias.Table[key] = alias.New(value)
		} else {
			return L.Push(nil, err)
		}
	case lua.LUA_TFUNCTION:
		chank := L.Dump()
		alias.Table[key] = &LuaBinaryChank{Chank: chank}
	}
	return L.Push(true)
}

func cmdGetAlias(L lua.Lua) int {
	name, nameErr := L.ToString(-1)
	if nameErr != nil {
		return L.Push(nil, nameErr)
	}
	value, ok := alias.Table[name]
	if !ok {
		L.PushNil()
		return 1
	}
	switch v := value.(type) {
	case *LuaBinaryChank:
		if err := L.LoadBufferX(name, v.Chank, "b"); err != nil {
			return L.Push(nil, err.Error())
		}
	default:
		L.PushString(v.String())
	}
	return 1
}

func cmdSetEnv(L lua.Lua) int {
	name, nameErr := L.ToString(-2)
	if nameErr != nil {
		return L.Push(nil, nameErr)
	}
	value, valueErr := L.ToString(-1)
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

func cmdGetEnv(L lua.Lua) int {
	name, nameErr := L.ToString(-1)
	if nameErr != nil {
		return L.Push(nil)
	}
	value, ok := interpreter.OurGetEnv(name)
	if ok && len(value) > 0 {
		L.PushString(value)
	} else {
		L.PushNil()
	}
	return 1
}

func cmdExec(L lua.Lua) int {
	errorlevel := interpreter.NOERROR
	var err error
	if L.IsTable(1) {
		L.Len(1)
		n, _ := L.ToInteger(-1)
		L.Pop(1)
		args := make([]string, 0, n+1)
		for i := 0; i <= n; i++ {
			L.RawGetI(-1, lua.Integer(i))
			arg1, err := L.ToString(-1)
			if err == nil && arg1 != "" {
				args = append(args, arg1)
			}
			L.Pop(1)
		}
		it := getRegInt(L)
		if it == nil {
			it = interpreter.New()
			it.Tag = NewNyagosLua()
			it.CloneHook = func(this *interpreter.Interpreter) error {
				this.Tag = NewNyagosLua()
				it.CloseHook = func(this *interpreter.Interpreter) {
					if L, ok := it.Tag.(lua.Lua); ok {
						L.Close()
					}
					it.Tag = nil
				}
				return nil
			}
		} else {
			it, err = it.Clone()
			if err != nil {
				return L.Push(nil, err)
			}
		}
		it.Args = args
		errorlevel, err = it.Spawnvp()
	} else {
		statement, statementErr := L.ToString(1)
		if statementErr != nil {
			return L.Push(nil, statementErr)
		}
		it := getRegInt(L)
		if it == nil {
			it = interpreter.New()
			it.Tag = L
		}
		errorlevel, err = it.Interpret(statement)
	}
	return L.Push(int(errorlevel), err)
}

type emptyWriter struct{}

func (e *emptyWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func cmdEval(L lua.Lua) int {
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
		it.Tag = L
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

func luaStackToSlice(L lua.Lua) []string {
	argc := L.GetTop()
	argv := make([]string, 0, argc)
	for i := 1; i <= argc; i++ {
		if L.IsTable(i) {
			L.Len(i)
			size, size_err := L.ToInteger(-1)
			L.Pop(1)
			if size_err == nil {
				// zero element
				L.RawGetI(i, 0)
				if s, err := L.ToString(-1); err == nil && s != "" {
					argv = append(argv, s)
				}
				L.Pop(1)
				// 1,2,3...
				for j := 1; j <= size; j++ {
					L.RawGetI(i, lua.Integer(j))
					if s, err := L.ToString(-1); err == nil {
						argv = append(argv, s)
					} else {
						argv = append(argv, "")
					}
					L.Pop(1)
				}
			}
		} else {
			s, s_err := L.ToString(i)
			if s_err != nil {
				s = ""
			}
			argv = append(argv, s)
		}
	}
	return argv
}

func cmdRawExec(L lua.Lua) int {
	argv := luaStackToSlice(L)
	cmd1 := exec.Command(argv[0], argv[1:]...)
	cmd1.Stdout = os.Stdout
	cmd1.Stderr = os.Stderr
	cmd1.Stdin = os.Stdin
	if it := getRegInt(L); it != nil {
		if it.Stdout != nil {
			cmd1.Stdout = it.Stdout
		}
		if it.Stderr != nil {
			cmd1.Stderr = it.Stderr
		}
		if it.Stdin != nil {
			cmd1.Stdin = it.Stdin
		}
	}
	err := cmd1.Run()
	errorlevel, errorlevelOk := interpreter.GetErrorLevel(cmd1.ProcessState)
	if !errorlevelOk {
		errorlevel = 255
	}
	if err != nil {
		fmt.Fprintln(cmd1.Stderr, err.Error())
		return L.Push(errorlevel, err.Error())
	} else {
		return L.Push(errorlevel)
	}
}

func cmdRawEval(L lua.Lua) int {
	argv := luaStackToSlice(L)
	cmd1 := exec.Command(argv[0], argv[1:]...)
	out, err := cmd1.Output()
	if err != nil {
		return L.Push(nil, err.Error())
	} else {
		L.PushAnsiString(out)
		return 1
	}
}

func cmdWrite(L lua.Lua) int {
	var out io.Writer = os.Stdout
	cmd := getRegInt(L)
	if cmd != nil && cmd.Stdout != nil {
		out = cmd.Stdout
	}
	return cmdWriteSub(L, out)
}

func cmdWriteErr(L lua.Lua) int {
	var out io.Writer = os.Stderr
	cmd := getRegInt(L)
	if cmd != nil && cmd.Stderr != nil {
		out = cmd.Stderr
	}
	return cmdWriteSub(L, out)
}

func cmdWriteSub(L lua.Lua, out io.Writer) int {
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

func cmdGetwd(L lua.Lua) int {
	wd, err := dos.Getwd()
	if err == nil {
		return L.Push(wd)
	} else {
		return L.Push(nil, err)
	}
}

func cmdWhich(L lua.Lua) int {
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

func cmdAtoU(L lua.Lua) int {
	str, err := dos.AtoU(L.ToAnsiString(1))
	if err == nil {
		L.PushString(str)
		return 1
	} else {
		return 0
	}
}

func cmdUtoA(L lua.Lua) int {
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

func cmdGlob(L lua.Lua) int {
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

func cmdGetHistory(this lua.Lua) int {
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

func cmdSetRuneWidth(this lua.Lua) int {
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

func cmdShellExecute(this lua.Lua) int {
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

func cmdStat(L lua.Lua) int {
	path, pathErr := L.ToString(1)
	if pathErr != nil {
		return L.Push(nil, pathErr)
	}
	var stat os.FileInfo
	var path_ string
	if len(path) > 0 && path[len(path)-1] == '\\' {
		path_ = dos.Join(path, ".")
	} else {
		path_ = path
	}
	statErr := dos.ForFiles(path_, func(f *dos.FileInfo) bool {
		stat = f
		return false
	})
	if statErr != nil {
		return L.Push(nil, statErr)
	}
	if stat == nil {
		return L.Push(nil, fmt.Errorf("%s: failed to stat", path))
	}
	L.NewTable()
	L.PushString(stat.Name())
	L.SetField(-2, "name")
	L.PushInteger(lua.Integer(stat.Size()))
	L.SetField(-2, "size")
	L.PushBool(stat.IsDir())
	L.SetField(-2, "isdir")
	t := stat.ModTime()
	L.NewTable()
	L.PushInteger(lua.Integer(t.Year()))
	L.SetField(-2, "year")
	L.PushInteger(lua.Integer(t.Month()))
	L.SetField(-2, "month")
	L.PushInteger(lua.Integer(t.Day()))
	L.SetField(-2, "day")
	L.PushInteger(lua.Integer(t.Hour()))
	L.SetField(-2, "hour")
	L.PushInteger(lua.Integer(t.Minute()))
	L.SetField(-2, "minute")
	L.PushInteger(lua.Integer(t.Second()))
	L.SetField(-2, "second")
	L.SetField(-2, "mtime")
	return 1
}

func cmdAccess(L lua.Lua) int {
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

func cmdPathJoin(L lua.Lua) int {
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

func cmdCommonPrefix(L lua.Lua) int {
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

func cmdGetKey(L lua.Lua) int {
	keycode, scancode, shiftstatus := conio.GetKey()
	L.PushInteger(lua.Integer(keycode))
	L.PushInteger(lua.Integer(scancode))
	L.PushInteger(lua.Integer(shiftstatus))
	return 3
}

func cmdGetViewWidth(L lua.Lua) int {
	width, height := conio.GetScreenBufferInfo().ViewSize()
	L.PushInteger(lua.Integer(width))
	L.PushInteger(lua.Integer(height))
	return 2
}

func cmdOpenFile(L lua.Lua) int {
	path, path_err := L.ToString(-2)
	if path_err != nil {
		return L.Push(nil, path_err.Error())
	}
	mode, mode_err := L.ToString(-1)
	if mode_err != nil {
		return L.Push(nil, mode_err.Error())
	}
	fd, fd_err := ansicfile.Open(path, mode)
	if fd_err != nil {
		return L.Push(nil, fd_err)
	}
	L.PushStream(fd)
	L.PushNil()
	return 2
}
