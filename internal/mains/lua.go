//go:build !vanilla
// +build !vanilla

package mains

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/yuin/gopher-lua"

	"github.com/nyaosorg/glua-ole"
	"github.com/nyaosorg/go-readline-ny"

	"github.com/nyaosorg/nyagos/internal/completion"
	"github.com/nyaosorg/nyagos/internal/frame"
	"github.com/nyaosorg/nyagos/internal/functions"
	"github.com/nyaosorg/nyagos/internal/history"
	"github.com/nyaosorg/nyagos/internal/shell"
)

// Lua is the alias for Lua's state type.
type Lua = *lua.LState

func makeVirtualTable(L Lua, getter, setter func(Lua) int) lua.LValue {
	table := L.NewTable()
	metaTable := L.NewTable()
	L.SetField(metaTable, "__index", L.NewFunction(getter))
	L.SetField(metaTable, "__newindex", L.NewFunction(setter))
	L.SetMetatable(table, metaTable)
	return table
}

var numberProperty = map[string]*int{
	"histsize": &history.MaxSaveHistory,
}

var stringProperty = map[string]*string{
	"antihistquot": &history.DisableMarks,
	"histchar":     &history.Mark,
	"quotation":    &readline.Delimiters,
	"version":      &frame.Version,
}

var boolProperty = map[string]*bool{
	"silentmode":        &frame.SilentMode,
	"completion_hidden": &completion.IncludeHidden,
	"completion_slash":  &completion.UseSlash,
}

var funcPropertySetter = map[string](func(*lua.LFunction)){
	"preexechook": func(f *lua.LFunction) {
		setExecHook(f, &shell.PreExecHook)
	},
	"postexechook": func(f *lua.LFunction) {
		setExecHook(f, &shell.PostExecHook)
	},
}

func nyagosGetter(L Lua) int {
	keyTmp, ok := L.Get(2).(lua.LString)
	if !ok {
		return lerror(L, "nyagos[]: too few arguments")
	}
	key := string(keyTmp)
	if ptr, ok := numberProperty[key]; ok {
		L.Push(lua.LNumber(*ptr))
	} else if ptr, ok := stringProperty[key]; ok {
		L.Push(lua.LString(*ptr))
	} else if ptr, ok := boolProperty[key]; ok {
		if *ptr {
			L.Push(lua.LTrue)
		} else {
			L.Push(lua.LFalse)
		}
	} else {
		L.Push(L.RawGet(L.Get(1).(*lua.LTable), keyTmp))
	}
	return 1
}

func nyagosSetter(L Lua) int {
	keyTmp, ok := L.Get(2).(lua.LString)
	if !ok {
		return lerror(L, "nyagos[]: too few arguments")
	}
	key := string(keyTmp)
	if ptr, ok := numberProperty[key]; ok {
		val, ok := L.Get(3).(lua.LNumber)
		if !ok {
			return lerror(L, "nyagos[]: val is not a number")
		}
		*ptr = int(val)
	} else if ptr, ok := stringProperty[key]; ok {
		val, ok := L.Get(3).(lua.LString)
		if !ok {
			return lerror(L, "nyagos[]: val is not a string")
		}
		*ptr = string(val)
	} else if ptr, ok := boolProperty[key]; ok {
		if L.Get(3) == lua.LTrue {
			*ptr = true
		} else if L.Get(3) == lua.LFalse {
			*ptr = false
		} else {
			return lerror(L, fmt.Sprintf("nyagos.%s: must be boolean", key))
		}
	} else if setter, ok := funcPropertySetter[key]; ok {
		if f, ok := L.Get(3).(*lua.LFunction); ok {
			setter(f)
		} else if L.Get(3) == lua.LNil {
			setter(nil)
		} else {
			return lerror(L, fmt.Sprintf("nyagos.%s: must be function", key))
		}
	} else {
		L.RawSet(L.Get(1).(*lua.LTable), L.Get(2), L.Get(3))
	}
	return 1
}

var isHookSetup = false

func lerror(L Lua, s string) int {
	L.Push(lua.LNil)
	L.Push(lua.LString(s))
	return 2
}

// NewLua sets up the lua instance with NYAOGS' environment.
func NewLua() (Lua, error) {
	L := lua.NewState(
		lua.Options{
			IncludeGoStackTrace: true,
			SkipOpenLibs:        true,
		})
	lua.OpenBase(L)
	lua.OpenChannel(L)
	lua.OpenCoroutine(L)
	lua.OpenDebug(L)
	// skip lua.OpenIo(L)
	lua.OpenMath(L)
	lua.OpenOs(L)
	lua.OpenPackage(L)
	lua.OpenString(L)
	lua.OpenTable(L)

	ioTable := openIo(L)
	L.SetGlobal("io", ioTable)

	nyagosTable := L.NewTable()

	envTable := makeVirtualTable(L,
		lua2param(functions.CmdGetEnv),
		lua2param(functions.CmdSetEnv))
	L.SetField(nyagosTable, "env", envTable)

	compTable := makeVirtualTable(L,
		complete4getter,
		complete4setter)
	L.SetField(nyagosTable, "complete_for", compTable)

	aliasTable := makeVirtualTable(L, cmdGetAlias, cmdSetAlias)
	L.SetField(nyagosTable, "alias", aliasTable)
	L.SetField(nyagosTable, "setalias", L.NewFunction(cmdSetAlias))
	L.SetField(nyagosTable, "getalias", L.NewFunction(cmdGetAlias))

	for name, function := range functions.Table {
		L.SetField(nyagosTable, name, L.NewFunction(lua2param(function)))
	}

	optionTable := makeVirtualTable(L, lua2param(functions.GetOption), lua2param(functions.SetOption))
	L.SetField(nyagosTable, "option", optionTable)

	L.SetField(nyagosTable, "lines", L.GetField(ioTable, "lines"))
	L.SetField(nyagosTable, "open", L.GetField(ioTable, "open"))
	L.SetField(nyagosTable, "loadfile", L.GetGlobal("loadfile"))

	keyTable := makeVirtualTable(L, lua2cmd(functions.CmdGetBindKey), cmdBindKey)
	L.SetField(nyagosTable, "key", keyTable)
	L.SetField(nyagosTable, "bindkey", L.NewFunction(cmdBindKey))
	L.SetField(nyagosTable, "exec", L.NewFunction(cmdExec))
	L.SetField(nyagosTable, "eval", L.NewFunction(cmdEval))
	L.SetField(nyagosTable, "prompt", L.NewFunction(lua2param(functions.Prompt)))
	L.SetField(nyagosTable, "create_object", L.NewFunction(ole.CreateObject))
	L.SetField(nyagosTable, "to_ole_integer", L.NewFunction(ole.ToOleInteger))
	L.SetField(nyagosTable, "goarch", lua.LString(runtime.GOARCH))
	L.SetField(nyagosTable, "goversion", lua.LString(runtime.Version()))
	L.SetField(nyagosTable, "goos", lua.LString(runtime.GOOS))
	L.SetField(nyagosTable, "version", lua.LString(frame.Version))
	L.SetField(nyagosTable, "pathseparator", lua.LString(string(os.PathSeparator)))
	L.SetField(nyagosTable, "pathlistseperator", lua.LString(string(os.PathListSeparator)))

	if exePath, err := os.Executable(); err == nil {
		L.SetField(nyagosTable, "exe", lua.LString(exePath))
	} else {
		println("gopherSh: NewLua: os.Executable() failed: " + err.Error())
	}

	historyMeta := L.NewTable()
	L.SetField(historyMeta, "__index", L.NewFunction(lua2param(functions.CmdGetHistory)))
	L.SetField(historyMeta, "__len", L.NewFunction(lua2param(functions.CmdLenHistory)))
	historyTable := L.NewTable()
	L.SetMetatable(historyTable, historyMeta)
	L.SetField(nyagosTable, "history", historyTable)

	metaTable := L.NewTable()
	L.SetField(metaTable, "__index", L.NewFunction(nyagosGetter))
	L.SetField(metaTable, "__newindex", L.NewFunction(nyagosSetter))
	L.SetMetatable(nyagosTable, metaTable)

	L.SetGlobal("nyagos", nyagosTable)

	shareTable := L.NewTable()
	L.SetGlobal("share", shareTable)

	setupUtf8Table(L)

	bit32Table := L.NewTable()
	L.SetField(bit32Table, "band", L.NewFunction(lua2param(functions.CmdBitAnd)))
	L.SetField(bit32Table, "bor", L.NewFunction(lua2param(functions.CmdBitOr)))
	L.SetField(bit32Table, "bxor", L.NewFunction(lua2param(functions.CmdBitXor)))
	L.SetGlobal("bit32", bit32Table)

	L.SetGlobal("print", L.NewFunction(lua2param(functions.CmdPrint)))

	if !isHookSetup {
		orgArgHook = shell.SetArgsHook(newArgHook)

		orgOnCommandNotFound = shell.OnCommandNotFound
		shell.OnCommandNotFound = (&_LuaCallBack{Lua: L}).onCommandNotFound
		isHookSetup = true
	}

	return L, nil
}

func lvalueToInterface(L Lua, valueTmp lua.LValue) interface{} {
	if valueTmp == lua.LNil {
		return nil
	} else if valueTmp == lua.LTrue {
		return true
	} else if valueTmp == lua.LFalse {
		return false
	}
	switch value := valueTmp.(type) {
	case lua.LString:
		return string(value)
	case lua.LNumber:
		return float64(value)
	case *lua.LUserData:
		return value.Value
	case *lua.LFunction:
		return value
	case *lua.LTable:
		table := make(map[interface{}]interface{})
		L.ForEach(value, func(keyTmp, valTmp lua.LValue) {
			key := lvalueToInterface(L, keyTmp)
			val := lvalueToInterface(L, valTmp)
			if f, ok := key.(float64); ok {
				table[int(f)] = val
			} else {
				table[key] = val
			}
		})
		return table
	default:
		println("lvalueToInterface: type not found")
		println(reflect.TypeOf(value).String())
		return nil
	}
}

func luaArgsToInterfaces(L Lua) []interface{} {
	end := L.GetTop()
	var param []interface{}
	if end > 0 {
		param = make([]interface{}, end)
		for i := 0; i < end; i++ {
			param[i] = lvalueToInterface(L, L.Get(i+1))
		}
	} else {
		param = []interface{}{}
	}
	return param
}

// ToLValueT is the type which can get lua.LValue
type ToLValueT interface {
	ToLValue(Lua) lua.LValue
}

func interfaceToLValue(L Lua, valueTmp interface{}) lua.LValue {
	if valueTmp == nil {
		return lua.LNil
	}
	switch value := valueTmp.(type) {
	case ToLValueT:
		return value.ToLValue(L)
	case string:
		return lua.LString(value)
	case error:
		return lua.LString(value.Error())
	case int:
		return lua.LNumber(value)
	case int16:
		return lua.LNumber(value)
	case int32:
		return lua.LNumber(value)
	case int64:
		return lua.LNumber(value)
	case uint:
		return lua.LNumber(value)
	case uint16:
		return lua.LNumber(value)
	case uint32:
		return lua.LNumber(value)
	case uint64:
		return lua.LNumber(value)
	case uintptr:
		return lua.LNumber(value)
	case float32:
		return lua.LNumber(value)
	case float64:
		return lua.LNumber(value)
	case time.Month:
		return lua.LNumber(value)
	case bool:
		if value {
			return lua.LTrue
		}
		return lua.LFalse
	case func([]interface{}) []interface{}:
		return L.NewFunction(lua2cmd(value))
	case func(*functions.Param) []interface{}:
		return L.NewFunction(lua2param(value))
	case reflect.Value:
		switch value.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return lua.LNumber(value.Int())
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return lua.LNumber(value.Uint())
		case reflect.Float32, reflect.Float64:
			return lua.LNumber(value.Float())
		case reflect.Bool:
			if value.Bool() {
				return lua.LTrue
			}
			return lua.LFalse
		case reflect.String:
			return lua.LString(value.String())
		case reflect.Interface:
			return interfaceToLValue(L, value.Interface())
		default:
			panic("not supporting type even in reflect value: " + value.Kind().String())
		}
	default:
		reflectValue := reflect.ValueOf(value)
		switch reflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			elem := reflectValue.Type().Elem()
			if elem.Kind() == reflect.Uint8 {
				buffer := make([]byte, reflectValue.Len())
				for i, end := 0, reflectValue.Len(); i < end; i++ {
					buffer[i] = byte(reflectValue.Index(i).Uint())
				}
				return lua.LString(string(buffer))
			}
			array1 := L.NewTable()
			for i, end := 0, reflectValue.Len(); i < end; i++ {
				val := reflectValue.Index(i)
				L.SetTable(array1,
					interfaceToLValue(L, i+1),
					interfaceToLValue(L, val))
			}
			return array1
		case reflect.Map:
			map1 := L.NewTable()
			for _, key := range reflectValue.MapKeys() {
				L.SetTable(map1,
					interfaceToLValue(L, key),
					interfaceToLValue(L, reflectValue.MapIndex(key)))
			}
			return map1
		default:
			println("interfaceToLValue: not support type")
			println(reflect.TypeOf(value).String())
			return nil
		}

	}
}

func pushInterfaces(L Lua, values []interface{}) {
	for _, value := range values {
		L.Push(interfaceToLValue(L, value))
	}
}

func lua2cmd(f func([]interface{}) []interface{}) func(Lua) int {
	return func(L Lua) int {
		param := luaArgsToInterfaces(L)
		result := f(param)
		pushInterfaces(L, result)
		return len(result)
	}
}

type shellKeyT struct{}

var shellKey shellKeyT

func getRegInt(L Lua) (context.Context, *shell.Shell) {
	ctx := getContext(L)
	if ctx == nil {
		println("getRegInt: could not find context in Lua instance")
		return context.Background(), nil
	}
	sh, ok := ctx.Value(shellKey).(*shell.Shell)
	if !ok {
		println("getRegInt: could not find shell in Lua instance")
		return ctx, nil
	}
	return ctx, sh
}

func lua2param(f func(*functions.Param) []interface{}) func(Lua) int {
	return func(L Lua) int {
		_, sh := getRegInt(L)
		param := &functions.Param{
			Args: luaArgsToInterfaces(L),
		}
		if sh != nil {
			param.In = sh.In()
			param.Out = sh.Out()
			param.Err = sh.Err()
		} else {
			param.In = os.Stdin
			param.Out = os.Stdout
			param.Err = os.Stderr
		}
		param.Term = colorable.NewColorableStdout()
		result := f(param)
		pushInterfaces(L, result)
		return len(result)
	}
}

const ctxkey = "github.com/nyaosorg/nyagos"

// setContext
func setContext(ctx context.Context, L Lua) {
	reg := L.Get(lua.RegistryIndex)
	if ctx != nil {
		u := L.NewUserData()
		u.Value = ctx
		L.SetField(reg, ctxkey, u)

		L.SetContext(ctx)
	} else {
		L.SetField(reg, ctxkey, lua.LNil)

		L.SetContext(context.Background())
	}
}

func getContext(L Lua) context.Context {
	reg := L.Get(lua.RegistryIndex)
	valueUD, ok := L.GetField(reg, ctxkey).(*lua.LUserData)
	if !ok {
		return nil
	}
	ctx, ok := valueUD.Value.(context.Context)
	if !ok {
		return nil
	}
	return ctx
}

func dispose(L *lua.LState, val lua.LValue) {
	gc := L.GetMetaField(val, "__gc")
	if f, ok := gc.(*lua.LFunction); ok {
		L.Push(f)
		L.Push(val)
		L.PCall(1, 0, nil)
	}
}

type _XFile struct {
	File      *os.File
	br        *bufio.Reader
	dontClose bool
	closed    bool
	eof       bool
}

func (xf *_XFile) reader() *bufio.Reader {
	if xf.br == nil {
		if xf.File != nil {
			xf.br = bufio.NewReader(xf.File)
		} else {
			panic("_XFile.reader() not found reader object")
		}
	}
	return xf.br
}

func (xf *_XFile) Write(b []byte) (int, error) { return xf.File.Write(b) }
func (xf *_XFile) Read(b []byte) (int, error)  { return xf.reader().Read(b) }
func (xf *_XFile) ReadByte() (byte, error)     { return xf.reader().ReadByte() }
func (xf *_XFile) UnreadByte() error           { return xf.reader().UnreadByte() }
func (xf *_XFile) ReadString(d byte) (string, error) {
	return xf.reader().ReadString(d)
}

func (xf *_XFile) Seek(offset int64, whence int) (int64, error) {
	xf.eof = false
	if xf.br != nil {
		back := xf.br.Buffered()
		if back > 0 {
			xf.File.Seek(int64(-back), io.SeekCurrent)
		}
		xf.br = nil
	}
	return xf.File.Seek(offset, whence)
}

func (xf *_XFile) Close() error {
	xf.br = nil
	if !xf.dontClose && !xf.closed {
		xf.closed = true
		return xf.File.Close()
	}
	return nil
}

func (xf *_XFile) IsClosed() bool { return xf.closed }
func (xf *_XFile) EOF() bool      { return xf.eof }
func (xf *_XFile) SetEOF()        { xf.eof = true }

func (xf *_XFile) Sync() error {
	return xf.File.Sync()
}

func luaRedirect(ctx context.Context, _stdin, _stdout, _stderr *os.File, L Lua, callback func() error) error {
	ioTbl := L.GetGlobal("io")

	orgStdin := L.GetField(ioTbl, "stdin")
	orgStdout := L.GetField(ioTbl, "stdout")
	orgStderr := L.GetField(ioTbl, "stderr")

	stdin := newXFile(L, &_XFile{File: _stdin, dontClose: true}, true, false)
	stdout := newXFile(L, &_XFile{File: _stdout, dontClose: true}, false, true)
	stderr := newXFile(L, &_XFile{File: _stderr, dontClose: true}, false, true)

	L.SetField(ioTbl, "stdin", stdin)
	L.SetField(ioTbl, "stdout", stdout)
	L.SetField(ioTbl, "stderr", stderr)

	err := callback()

	dispose(L, stdin)
	dispose(L, stdout)
	dispose(L, stderr)
	L.SetField(ioTbl, "stdin", orgStdin)
	L.SetField(ioTbl, "stdout", orgStdout)
	L.SetField(ioTbl, "stderr", orgStderr)
	return err
}

func execLuaKeepContextAndShell(ctx context.Context, sh *shell.Shell, L Lua, nargs, nresult int) error {
	defer setContext(getContext(L), L)
	ctx = context.WithValue(ctx, shellKey, sh)
	setContext(ctx, L)

	return luaRedirect(ctx, sh.Stdio[0], sh.Stdio[1], sh.Stdio[2], L, func() error {
		return L.PCall(nargs, nresult, nil)
	})
}
