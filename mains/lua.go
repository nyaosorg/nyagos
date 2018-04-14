package mains

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"

	"github.com/zetamatta/nyagos/commands"
	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/functions"
	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/lua"
	ole "github.com/zetamatta/nyagos/lua/ole"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

type shellKeyT struct{}

var shellKey shellKeyT

func getRegInt(L lua.Lua) (context.Context, *shell.Shell) {
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

func callLua(ctx context.Context, sh *shell.Shell, nargs int, nresult int) error {
	luawrapper, ok := sh.Tag().(*luaWrapper)
	if !ok {
		return errors.New("callLua: can not find Lua instance in the shell")
	}
	L := luawrapper.Lua
	ctx = context.WithValue(ctx, shellKey, sh)
	return L.CallWithContext(ctx, nargs, nresult)
}

var orgArgHook func(context.Context, *shell.Shell, []string) ([]string, error)

var luaArgsFilter lua.Object = lua.TNil{}

func newArgHook(ctx context.Context, it *shell.Shell, args []string) ([]string, error) {
	luawrapper, ok := it.Tag().(*luaWrapper)
	if !ok {
		return nil, errors.New("Could not get lua instance(newArgHook)")
	}
	L := luawrapper.Lua
	L.Push(luaArgsFilter)
	if !L.IsFunction(-1) {
		L.Pop(1)
		return orgArgHook(ctx, it, args)
	}
	L.NewTable()
	for i := 0; i < len(args); i++ {
		L.PushString(args[i])
		L.RawSetI(-2, lua.Integer(i))
	}
	if err := callLua(ctx, it, 1, 1); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return orgArgHook(ctx, it, args)
	}
	if L.GetType(-1) != lua.LUA_TTABLE {
		return orgArgHook(ctx, it, args)
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
	return orgArgHook(ctx, it, newargs)
}

var orgOnCommandNotFound func(*shell.Cmd, error) error

var luaOnCommandNotFound lua.Object = lua.TNil{}

func on_command_not_found(inte *shell.Cmd, err error) error {
	luawrapper, ok := inte.Tag().(*luaWrapper)
	if !ok {
		return errors.New("Could get lua instance(on_command_not_found)")
	}
	L := luawrapper.Lua

	L.Push(luaOnCommandNotFound)
	if !L.IsFunction(-1) {
		L.Pop(1)
		return orgOnCommandNotFound(inte, err)
	}
	L.NewTable()
	for key, val := range inte.Args() {
		L.PushString(val)
		L.RawSetI(-2, lua.Integer(key))
	}
	err1 := L.Call(1, 1)
	defer L.Pop(1)
	if err1 != nil {
		return err
	}
	if L.ToBool(-1) {
		return nil
	} else {
		return orgOnCommandNotFound(inte, err)
	}
}

func getOption(L lua.Lua) int {
	key, key_err := L.ToString(2)
	if key_err != nil {
		return L.Push(nil, key_err)
	}
	ptr, ok := commands.BoolOptions[key]
	if !ok {
		return L.Push(nil)
	}
	return L.Push(&lua.BoolProperty{Pointer: ptr})
}

func setOption(L lua.Lua) int {
	key, key_err := L.ToString(2)
	if key_err != nil {
		return L.Push(nil, key_err)
	}
	ptr, ok := commands.BoolOptions[key]
	if !ok {
		return L.Push(nil)
	}
	opt := lua.BoolProperty{Pointer: ptr}
	if err := opt.Set(L, 3); err != nil {
		return L.Push(nil, err.Error())
	} else {
		return L.Push(true)
	}
}

var nyagos_table_member map[string]lua.Object

func getNyagosTable(L lua.Lua) int {
	index, index_err := L.ToString(2)
	if index_err != nil {
		return L.Push(nil, index_err.Error())
	}
	if entry, entry_ok := nyagos_table_member[index]; entry_ok {
		return L.Push(entry)
	} else if index == "exe" {
		if exeName, exeNameErr := os.Executable(); exeNameErr != nil {
			return L.Push(nil, exeNameErr.Error())
		} else {
			L.PushString(exeName)
			return 1
		}
	} else {
		L.PushNil()
		return 1
	}
}

type IProperty interface {
	Push(lua.Lua) int
	Set(lua.Lua, int) error
}

func setNyagosTable(L lua.Lua) int {
	index, index_err := L.ToString(2)
	if index_err != nil {
		return L.Push(nil, index_err)
	}
	if current_value, exists := nyagos_table_member[index]; exists {
		if property, castOk := current_value.(IProperty); castOk {
			if err := property.Set(L, 3); err != nil {
				fmt.Fprintf(os.Stderr, "nyagos.%s: %s\n", index, err.Error())
				return L.Push(nil, err)
			} else {
				return L.Push(true)
			}
		} else {
			value, value_err := L.ToObject(3)
			if value_err != nil {
				return L.Push(nil, value_err)
			}
			nyagos_table_member[index] = value
			return L.Push(true)
		}
	} else {
		fmt.Fprintf(os.Stderr, "nyagos.%s: reserved variable.\n", index)
		return L.Push(nil)
	}
}

var share_table = map[string]lua.Object{}
var share_table_generation = map[string]int{}

func setMemberOfShareTable(L lua.Lua) int {
	// table exists at [-3]
	key, err := L.ToObject(-2)
	if err != nil {
		return L.Push(nil, err)
	}
	val, err := L.ToObject(-1)
	if err != nil {
		return L.Push(nil, err)
	}
	L.RawSet(-3) // pop 2
	L.GetMetaTable(-1)
	L.GetField(-1, "..")
	parentkey, err := L.ToString(-1)
	if err != nil {
		println(err.Error())
		return L.Push(nil, err)
	}
	L.Pop(1) // drop string
	L.GetField(-1, "age")
	age, err := L.ToInteger(-1)
	L.Pop(2) // drop integer and metatable

	if err != nil || age != share_table_generation[parentkey] {
		// println("old variable")
		return 0
	}

	table1, ok := share_table[parentkey]
	if !ok {
		err := fmt.Errorf("%s: not found in share_table()", parentkey)
		println(err.Error())
		return L.Push(nil, err.Error())
	}
	if t, ok := table1.(*lua.MetaTableOwner); ok {
		table1 = t.Body
	}
	table2, ok := table1.(*lua.TTable)
	if !ok {
		err := fmt.Errorf("%s: not table in share_table()", parentkey)
		type1 := reflect.TypeOf(table1)
		println(type1.String())
		println(err.Error())
		return L.Push(nil, err.Error())
	}
	switch t := key.(type) {
	case lua.TString:
		table2.Dict[string(t)] = val
	case lua.TRawString:
		table2.Dict[string(t)] = val
	case lua.Integer:
		table2.Array[int(t)] = val
	}
	return 0
}

func getShareTable(L lua.Lua) int {
	key, keyErr := L.ToString(-1)
	if keyErr != nil {
		return L.Push(nil, keyErr)
	}
	if value, ok := share_table[key]; ok {
		L.Push(value)
		if L.IsTable(-1) {
			L.NewTable()
			L.PushGoFunction(setMemberOfShareTable)
			L.SetField(-2, "__newindex")
			L.PushString(key)
			L.SetField(-2, "..")
			L.PushInteger(lua.Integer(share_table_generation[key]))
			L.SetField(-2, "age")
			L.SetMetaTable(-2)
		}
		return 1

	} else {
		L.PushNil()
		return 1
	}
}

func setShareTable(L lua.Lua) int {
	key, keyErr := L.ToString(-2)
	if keyErr != nil {
		return L.Push(nil, keyErr)
	}
	value, valErr := L.ToObject(-1)
	if valErr != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", key, valErr.Error())
		return L.Push(nil, valErr)
	}
	share_table[key] = value
	share_table_generation[key]++
	return 1
}

var hook_setuped = false

func NewLua() (lua.Lua, error) {
	this, err := lua.New()
	if err != nil {
		return this, err
	}
	this.OpenLibs()

	this.Push(&lua.VirtualTable{
		Name:     "nyagos",
		Index:    getNyagosTable,
		NewIndex: setNyagosTable})
	this.SetGlobal("nyagos")

	this.Push(&lua.VirtualTable{
		Name:     "share",
		Index:    getShareTable,
		NewIndex: setShareTable})
	this.SetGlobal("share")

	this.PushGoFunction(lua2param(functions.CmdPrint))
	this.SetGlobal("print")

	if !hook_setuped {
		orgArgHook = shell.SetArgsHook(newArgHook)

		orgOnCommandNotFound = shell.OnCommandNotFound
		shell.OnCommandNotFound = on_command_not_found
		hook_setuped = true
	}
	return this, nil
}

func luaArgsToInterfaces(L lua.Lua) []interface{} {
	end := L.GetTop()
	var param []interface{}
	if end > 0 {
		param = make([]interface{}, 0, end-1)
		for i := 1; i <= end; i++ {
			value, _ := L.ToInterface(i)
			param = append(param, value)
		}
	} else {
		param = []interface{}{}
	}
	return param
}

func pushInterfaces(L lua.Lua, values []interface{}) {
	for _, value := range values {
		L.PushReflect(value)
	}
}

func lua2cmd(f func([]interface{}) []interface{}) func(lua.Lua) int {
	return func(L lua.Lua) int {
		param := luaArgsToInterfaces(L)
		result := f(param)
		pushInterfaces(L, result)
		return len(result)
	}
}

func lua2param(f func(*functions.Param) []interface{}) func(lua.Lua) int {
	return func(L lua.Lua) int {
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

func init() {
	nyagos_table_member = map[string]lua.Object{
		"alias": &lua.VirtualTable{
			Name:     "nyagos.alias",
			Index:    cmdGetAlias,
			NewIndex: cmdSetAlias},
		"antihistquot": lua.StringProperty{Pointer: &history.DisableMarks},
		"argsfilter":   lua.Property{Pointer: &luaArgsFilter},
		"key": &lua.VirtualTable{
			Name:     "nyagos.key",
			Index:    cmdGetBindKey,
			NewIndex: cmdBindKey},
		"bindkey":           lua.TGoFunction(cmdBindKey),
		"completion_slash":  lua.BoolProperty{Pointer: &completion.UseSlash},
		"completion_hook":   lua.Property{Pointer: &completionHook},
		"completion_hidden": lua.BoolProperty{Pointer: &completion.IncludeHidden},
		"create_object":     lua.TGoFunction(ole.CreateObject),
		"env": &lua.VirtualTable{
			Name:     "nyagos.env",
			Index:    lua2cmd(functions.CmdGetEnv),
			NewIndex: lua2cmd(functions.CmdSetEnv)},
		"eval":      lua.TGoFunction(cmdEval),
		"exec":      lua.TGoFunction(cmdExec),
		"filter":    lua.Property{Pointer: &luaFilter},
		"getalias":  lua.TGoFunction(cmdGetAlias),
		"goarch":    lua.TString(runtime.GOARCH),
		"goversion": lua.TString(runtime.Version()),
		"histchar":  lua.StringProperty{Pointer: &history.Mark},
		"history": &lua.VirtualTable{
			Name:  "nyagos.history",
			Index: lua2cmd(functions.CmdGetHistory),
			Len:   lua2cmd(functions.CmdLenHistory)},
		"lines":                lua.TGoFunction(cmdLines),
		"loadfile":             lua.TGoFunction(cmdLoadFile),
		"on_command_not_found": lua.Property{Pointer: &luaOnCommandNotFound},
		"open":                 lua.TGoFunction(cmdOpenFile),
		"option": &lua.VirtualTable{
			Name:     "nyagos.option",
			Index:    getOption,
			NewIndex: setOption},
		"prompt":     lua.Property{Pointer: &prompt_hook},
		"quotation":  lua.StringProperty{Pointer: &readline.Delimiters},
		"setalias":   lua.TGoFunction(cmdSetAlias),
		"silentmode": &lua.BoolProperty{Pointer: &frame.SilentMode},
		"version":    lua.StringProperty{Pointer: &frame.Version},
	}
	for key, val := range functions.Table {
		nyagos_table_member[key] = lua.TGoFunction(lua2cmd(val))
	}
	for key, val := range functions.Table2 {
		nyagos_table_member[key] = lua.TGoFunction(lua2param(val))
	}
}

func runLua(ctx context.Context, it *shell.Shell, L lua.Lua, fname string) ([]byte, error) {
	_, err := os.Stat(fname)
	if err != nil {
		if os.IsNotExist(err) {
			// println("pass " + fname + " (not exists)")
			return []byte{}, nil
		} else {
			return nil, err
		}
	}
	if _, err := L.LoadFile(fname, "bt"); err != nil {
		return nil, err
	}
	chank := L.Dump()
	if err := callLua(ctx, it, 0, 0); err != nil {
		return nil, err
	}
	// println("Run: " + fname)
	return chank, nil
}
