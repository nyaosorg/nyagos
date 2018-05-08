package gopherSh

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"time"

	"github.com/yuin/gopher-lua"

	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/functions"
	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

type Lua = *lua.LState

func makeVirtualTable(L Lua, getter, setter func(Lua) int) lua.LValue {
	table := L.NewTable()
	metaTable := L.NewTable()
	L.SetField(metaTable, "__index", L.NewFunction(getter))
	L.SetField(metaTable, "__newindex", L.NewFunction(setter))
	L.SetMetatable(table, metaTable)
	return table
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

func nyagosGetter(L Lua) int {
	keyTmp, ok := L.Get(2).(lua.LString)
	if !ok {
		return lerror(L, "nyagos[]: too few arguments")
	}
	key := string(keyTmp)
	if ptr, ok := stringProperty[key]; ok {
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
	if ptr, ok := stringProperty[key]; ok {
		val, ok := L.Get(3).(lua.LString)
		if !ok {
			return lerror(L, fmt.Sprintf("nyagos[]: val is not a string"))
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
	} else {
		L.RawSet(L.Get(1).(*lua.LTable), L.Get(2), L.Get(3))
	}
	return 1
}

var isHookSetup = false

func NewLua() (Lua, error) {
	L := lua.NewState(lua.Options{IncludeGoStackTrace: true})

	nyagosTable := L.NewTable()

	for name, function := range functions.Table {
		L.SetField(nyagosTable, name, L.NewFunction(lua2cmd(function)))
	}
	envTable := makeVirtualTable(L,
		lua2cmd(functions.CmdGetEnv),
		lua2cmd(functions.CmdSetEnv))
	L.SetField(nyagosTable, "env", envTable)

	aliasTable := makeVirtualTable(L, cmdGetAlias, cmdSetAlias)
	L.SetField(nyagosTable, "alias", aliasTable)
	L.SetField(nyagosTable, "setalias", L.NewFunction(cmdSetAlias))
	L.SetField(nyagosTable, "getalias", L.NewFunction(cmdGetAlias))

	for name, function := range functions.Table2 {
		L.SetField(nyagosTable, name, L.NewFunction(lua2param(function)))
	}

	optionTable := makeVirtualTable(L, lua2cmd(functions.GetOption), lua2cmd(functions.SetOption))
	L.SetField(nyagosTable, "option", optionTable)

	ioTable := L.GetGlobal("io")
	L.SetField(nyagosTable, "open", L.GetField(ioTable, "open"))
	L.SetField(nyagosTable, "lines", L.GetField(ioTable, "lines"))
	L.SetField(nyagosTable, "loadfile", L.GetGlobal("loadfile"))

	keyTable := makeVirtualTable(L, lua2cmd(functions.CmdGetBindKey), cmdBindKey)
	L.SetField(nyagosTable, "key", keyTable)
	L.SetField(nyagosTable, "bindkey", L.NewFunction(cmdBindKey))
	L.SetField(nyagosTable, "exec", L.NewFunction(cmdExec))
	L.SetField(nyagosTable, "eval", L.NewFunction(cmdEval))
	L.SetField(nyagosTable, "prompt", L.NewFunction(lua2cmd(functions.Prompt)))
	L.SetField(nyagosTable, "create_object", L.NewFunction(CreateObject))
	L.SetField(nyagosTable, "goarch", lua.LString(runtime.GOARCH))
	L.SetField(nyagosTable, "goversion", lua.LString(runtime.Version()))
	L.SetField(nyagosTable, "version", lua.LString(frame.Version))

	if exePath, err := os.Executable(); err == nil {
		L.SetField(nyagosTable, "exe", lua.LString(exePath))
	} else {
		println("gopherSh: NewLua: os.Executable() failed: " + err.Error())
	}

	historyMeta := L.NewTable()
	L.SetField(historyMeta, "__index", L.NewFunction(lua2cmd(functions.CmdGetHistory)))
	L.SetField(historyMeta, "__len", L.NewFunction(lua2cmd(functions.CmdLenHistory)))
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

	L.SetGlobal("print", L.NewFunction(lua2param(functions.CmdPrint)))

	if !isHookSetup {
		orgArgHook = shell.SetArgsHook(newArgHook)

		orgOnCommandNotFound = shell.OnCommandNotFound
		shell.OnCommandNotFound = onCommandNotFound
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
		return int(value)
	case *lua.LUserData:
		return value.Value
	case *lua.LFunction:
		return value
	case *lua.LTable:
		table := make(map[interface{}]interface{})
		L.ForEach(value, func(keyTmp, valTmp lua.LValue) {
			key := lvalueToInterface(L, keyTmp)
			val := lvalueToInterface(L, valTmp)
			table[key] = val
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
		param = make([]interface{}, 0, end-1)
		for i := 1; i <= end; i++ {
			param = append(param, lvalueToInterface(L, L.Get(i)))
		}
	} else {
		param = []interface{}{}
	}
	return param
}

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
	case int64:
		return lua.LNumber(value)
	case time.Month:
		return lua.LNumber(value)
	case bool:
		if value {
			return lua.LTrue
		} else {
			return lua.LFalse
		}
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
		case reflect.Bool:
			if value.Bool() {
				return lua.LTrue
			} else {
				return lua.LFalse
			}
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
				buffer := make([]byte, 0, reflectValue.Len())
				for i, end := 0, reflectValue.Len(); i < end; i++ {
					buffer = append(buffer, byte(reflectValue.Index(i).Uint()))
				}
				return lua.LString(string(buffer))
			} else {
				array1 := L.NewTable()
				for i, end := 0, reflectValue.Len(); i < end; i++ {
					val := reflectValue.Index(i)
					L.SetTable(array1,
						interfaceToLValue(L, i+1),
						interfaceToLValue(L, val))
				}
				return array1
			}
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
		result := f(param)
		pushInterfaces(L, result)
		return len(result)
	}
}

const ctxkey = "github.com/zetamatta/nyagos"

// setContext
// We does not use (lua.LState)SetContext.
// Because sometimes cancel is requrested on unexpected timing.
func setContext(L Lua, ctx context.Context) {
	reg := L.Get(lua.RegistryIndex)
	if ctx != nil {
		u := L.NewUserData()
		u.Value = ctx
		L.SetField(reg, ctxkey, u)
	} else {
		L.SetField(reg, ctxkey, lua.LNil)
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

func callCSL(ctx context.Context, sh *shell.Shell, L Lua, nargs, nresult int) (err error) {
	defer setContext(L, getContext(L))
	ctx = context.WithValue(ctx, shellKey, sh)
	setContext(L, ctx)
	return L.PCall(nargs, nresult, nil)
}

func callLua(ctx context.Context, sh *shell.Shell, nargs, nresult int) error {
	luawrapper, ok := sh.Tag().(*luaWrapper)
	if !ok {
		return errors.New("callLua: can not find Lua instance in the shell")
	}
	return callCSL(ctx, sh, luawrapper.Lua, nargs, nresult)
}
