package mains

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode"

	"github.com/mattn/go-colorable"
	"github.com/zetamatta/go-ansicfile"

	"github.com/zetamatta/nyagos/alias"
	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/lua"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

type LuaBinaryChank struct {
	Chank []byte
}

func (this *LuaBinaryChank) String() string {
	return "(lua-function)"
}

func (this *LuaBinaryChank) Call(ctx context.Context, cmd *shell.Cmd) (int, error) {
	L, L_ok := cmd.Tag.(lua.Lua)
	if !L_ok {
		return 255, errors.New("LuaBinaryChank.Call: Lua instance not found")
	}

	if f := cmd.Stdout; f != os.Stdout && f != os.Stderr {
		L.GetGlobal("io")        // +1
		L.GetField(-1, "output") // +1 (get function pointer)
		if err := L.PushFileWriter(f); err != nil {
			L.Pop(2)
			return 255, err
		}
		L.Call(1, 0)
		L.Pop(1) // remove io-table
	}
	if f := cmd.Stdin; f != os.Stdin {
		L.GetGlobal("io")       // +1
		L.GetField(-1, "input") // +1 (get function pointer)
		if err := L.PushFileReader(f); err != nil {
			L.Pop(2)
			return 255, err
		}
		L.Call(1, 0)
		L.Pop(1) // remove io-table
	}

	if err := L.LoadBufferX(cmd.Args[0], this.Chank, "b"); err != nil {
		return 255, err
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
	errorlevel := 0
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
				errorlevel = 255
				err = err1
			} else {
				it.Args = newargs
				errorlevel, err = it.SpawnvpContext(ctx)
			}
		} else if val, err1 := L.ToInteger(-1); err1 == nil {
			errorlevel = val
		} else if val, err1 := L.ToString(-1); val != "" && err1 == nil {
			it, err1 := cmd.Clone()
			if err1 != nil {
				errorlevel = 255
				err = err1
			} else {
				errorlevel, err = it.InterpretContext(ctx, val)
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

func cmdExec(L lua.Lua) int {
	errorlevel := 0
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
			println("main/lua_cmd.go: cmdExec: not found interpreter object")
			it = shell.New()
			it.Tag = L
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
			it = shell.New()
			it.Tag = L
		}
		errorlevel, err = it.Interpret(statement)
	}
	return L.Push(int(errorlevel), err)
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
		it := shell.New()
		it.Tag = L
		it.Stdout = w
		it.Interpret(statement)
		w.Close()
	}(statement, w)

	result, err := ioutil.ReadAll(r)
	r.Close()
	if err == nil {
		L.PushBytes(bytes.Trim(result, "\r\n\t "))
	} else {
		L.PushNil()
	}
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
	errorlevel, errorlevelOk := dos.GetErrorLevel(cmd1)
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
		L.PushBytes(out)
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
	switch f := out.(type) {
	case *os.File:
		out = colorable.NewColorable(f)
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

func cmdGetHistory(this lua.Lua) int {
	if default_history == nil {
		return 0
	}
	if this.GetType(-1) == lua.LUA_TNUMBER {
		val, err := this.ToInteger(-1)
		if err != nil {
			return this.Push(nil, err.Error())
		}
		this.PushString(default_history.At(val))
	} else {
		this.PushInteger(lua.Integer(default_history.Len()))
	}
	return 1
}

func cmdLenHistory(this lua.Lua) int {
	if default_history == nil {
		return 0
	}
	this.PushInteger(lua.Integer(default_history.Len()))
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
	readline.SetCharWidth(rune(char), width)
	this.PushBool(true)
	return 1
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

func cmdOpenFile(L lua.Lua) int {
	path, path_err := L.ToString(1)
	if path_err != nil {
		return L.Push(nil, path_err.Error())
	}
	mode, mode_err := L.ToString(2)
	if mode_err != nil {
		return L.Push(nil, mode_err.Error())
	}
	if mode == "" {
		mode = "r"
	}
	fd, fd_err := ansicfile.Open(path, mode)
	if fd_err != nil {
		// print("-- open '", path, "' failed --\n")
		return L.Push(nil, fd_err)
	}
	// print("-- open '", path, "' success --\n")
	L.PushStream(fd)
	return 1
}

func cmdLoadFile(L lua.Lua) int {
	path, path_err := L.ToString(-1)
	if path_err != nil {
		return L.Push(nil, path_err.Error())
	}
	_, err := L.LoadFile(path, "bt")
	if err != nil {
		return L.Push(nil, path+": "+err.Error())
	} else {
		return 1
	}
}

type iolines_t struct {
	Fd         *os.File
	Reader     *bufio.Reader
	Marks      []string
	HasToClose bool
}

func (this *iolines_t) Close() {
	this.Reader = nil
	if this.HasToClose {
		this.Fd.Close()
	}
	this.Fd = nil
}

func (this *iolines_t) Ok() bool {
	return this != nil && this.Fd != nil && this.Reader != nil
}

func iolines_t_gc(L lua.Lua) int {
	defer L.DeleteUserDataAnchor(1)
	userdata := iolines_t{}
	sync := L.ToUserDataTo(1, &userdata)
	defer sync()
	// print("iolines_t_gc: gc\n")
	if !userdata.Ok() {
		// print("iolines_t_gc: nil\n")
		return 0
	}
	userdata.Close()
	return 0
}

func cmdLinesCallback(L lua.Lua) int {
	userdata := iolines_t{}
	sync := L.ToUserDataTo(1, &userdata)
	defer sync()

	if !userdata.Ok() {
		// print("cmdLinesCallback: nil\n")
		return L.Push(nil)
	}
	count := 0
	var err error

	for _, mark := range userdata.Marks {
		if err != nil {
			break
		}
		if unicode.IsDigit(rune(mark[0])) {
			var nbytes int
			nbytes, err = strconv.Atoi(mark)
			if err != nil {
				break
			}
			line := make([]byte, nbytes)
			var nreads int
			nreads, err = userdata.Reader.Read(line)
			if nreads > 0 {
				if nreads < nbytes {
					L.PushBytes(line[:nreads])
				} else {
					L.PushBytes(line)
				}
				count++
			}
			break
		}
		switch mark {
		default:
			var line []byte
			line, err = userdata.Reader.ReadBytes('\n')
			if err != nil {
				break
			}
			if mark != "L" {
				line = bytes.TrimRight(line, "\r\n")
			}
			L.PushBytes(line)
			count++
			break
		case "a":
			var data []byte
			data, err = ioutil.ReadAll(userdata.Reader)
			if err != nil || len(data) <= 0 {
				userdata.Close()
				return L.Push(nil)
			}
			L.PushBytes(data)
			count++
			break
		case "n":
			var val int
			if _, err := fmt.Fscan(userdata.Reader, &val); err != nil {
				userdata.Close()
				return L.Push(nil)
			} else {
				L.PushInteger(lua.Integer(val))
				count++
				break
			}
		}
	}
	if err != nil {
		// print("cmdLinesCallback: eof\n")
		userdata.Close()
		return L.Push(nil)
	}
	// print("cmdLinesCallback: text='", text, "'\n")
	return count
}

func cmdLines(L lua.Lua) int {
	top := L.GetTop()
	if top < 1 || L.IsNil(1) {
		L.Push(cmdLinesCallback)
		cmd := getRegInt(L)
		L.PushUserData(&iolines_t{
			Fd:         cmd.Stdin,
			Reader:     bufio.NewReader(cmd.Stdin),
			HasToClose: false,
			Marks:      []string{"l"},
		})
		L.NewTable()
		L.Push(iolines_t_gc)
		L.SetField(-2, "__gc")
		L.SetMetaTable(-2)
		return 2
	}
	path, path_err := L.ToString(1)
	if path_err != nil {
		return L.Push(nil, path_err.Error())
	}
	var marks []string
	if top < 2 {
		marks = []string{"l"}
	} else {
		marks = make([]string, 0, top-2)
		for i := 2; i <= top; i++ {
			mark1, err1 := L.ToString(i)
			if err1 != nil {
				return L.Push(nil, err1.Error())
			}
			marks = append(marks, mark1)
		}
	}
	fd, fd_err := os.Open(path)
	if fd_err != nil {
		return L.Push(nil, fd_err.Error())
	}
	L.Push(cmdLinesCallback)
	userdata := iolines_t{
		Fd:         fd,
		Reader:     bufio.NewReader(fd),
		Marks:      marks,
		HasToClose: true,
	}
	L.PushUserData(&userdata)
	L.NewTable()
	L.Push(iolines_t_gc)
	L.SetField(-2, "__gc")
	L.SetMetaTable(-2)
	// print("cmdLines: end\n")
	return 2
}
