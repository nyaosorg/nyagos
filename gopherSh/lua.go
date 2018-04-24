package main

import (
	"context"
	"errors"
	"os"

	"github.com/yuin/gopher-lua"

	"github.com/zetamatta/nyagos/functions"
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

func NewLua() (Lua, error) {
	L := lua.NewState()

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

	for name, function := range functions.Table2 {
		L.SetField(nyagosTable, name, L.NewFunction(lua2param(function)))
	}

	optionTable := makeVirtualTable(L, lua2cmd(functions.GetOption), lua2cmd(functions.SetOption))
	L.SetField(nyagosTable, "option", optionTable)

	ioTable := L.GetGlobal("io")
	L.SetField(nyagosTable, "open", L.GetField(ioTable, "open"))

	keyTable := makeVirtualTable(L, lua2cmd(functions.CmdGetBindKey), cmdBindKey)
	L.SetField(nyagosTable, "key", keyTable)
	L.SetField(nyagosTable, "bindkey", L.NewFunction(cmdBindKey))

	L.SetGlobal("nyagos", nyagosTable)

	shareTable := L.NewTable()
	L.SetGlobal("share", shareTable)

	L.SetGlobal("print", L.NewFunction(lua2param(functions.CmdPrint)))

	return L, nil
}

func luaArgsToInterfaces(L Lua) []interface{} {
	end := L.GetTop()
	var param []interface{}
	if end > 0 {
		param = make([]interface{}, 0, end-1)
		for i := 1; i <= end; i++ {
			switch L.Get(i).Type() {
			case lua.LTString:
				param = append(param, L.ToString(i))
			case lua.LTNumber:
				param = append(param, L.ToInt(i))
			case lua.LTBool:
				param = append(param, L.ToBool(i))
			case lua.LTFunction:
				param = append(param, L.ToFunction(i))
			case lua.LTUserData:
				param = append(param, L.ToUserData(i))
			case lua.LTNil:
				param = append(param, nil)
			default:
				param = append(param, nil)
			}
		}
	} else {
		param = []interface{}{}
	}
	return param
}

func pushInterfaces(L Lua, values []interface{}) {
	for _, valueTmp := range values {
		if valueTmp == nil {
			L.Push(lua.LNil)
		} else {
			switch value := valueTmp.(type) {
			case string:
				L.Push(lua.LString(value))
			case int:
				L.Push(lua.LNumber(value))
			case bool:
				if value {
					L.Push(lua.LTrue)
				} else {
					L.Push(lua.LFalse)
				}
			case func([]interface{}) []interface{}:
				L.Push(L.NewFunction(lua2cmd(value)))
			case func(*functions.Param) []interface{}:
				L.Push(L.NewFunction(lua2param(value)))
			}
		}
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
	ctx := L.Context()
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

func callCSL(ctx context.Context, sh *shell.Shell, L Lua, nargs, nresult int) (err error) {
	if save := L.Context(); save != nil {
		defer L.SetContext(save)
	}
	ctx = context.WithValue(ctx, shellKey, sh)
	L.SetContext(ctx)
	return L.PCall(nargs, nresult, nil)
}

func callLua(ctx context.Context, sh *shell.Shell, nargs, nresult int) error {
	luawrapper, ok := sh.Tag().(*luaWrapper)
	if !ok {
		return errors.New("callLua: can not find Lua instance in the shell")
	}
	return callCSL(ctx, sh, luawrapper.Lua, nargs, nresult)
}
