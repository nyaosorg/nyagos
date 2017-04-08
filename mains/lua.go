package mains

import (
	"fmt"
	"os"
	"runtime"
	"unsafe"

	"../history"
	"../interpreter"
	"../lua"
	ole "../lua/ole"
	"../readline"
)

const REGKEY_INTERPRETER = "nyagos.interpreter"

func setRegInt(L lua.Lua, it *interpreter.Interpreter) {
	L.PushValue(lua.LUA_REGISTRYINDEX)
	L.PushLightUserData(unsafe.Pointer(it))
	L.SetField(-2, REGKEY_INTERPRETER)
	L.Pop(1)
}

func getRegInt(L lua.Lua) *interpreter.Interpreter {
	L.PushValue(lua.LUA_REGISTRYINDEX)
	L.GetField(-1, REGKEY_INTERPRETER)
	rc := (*interpreter.Interpreter)(L.ToUserData(-1))
	L.Pop(2)
	return rc
}

func NyagosCallLua(L lua.Lua, it *interpreter.Interpreter, nargs int, nresult int) error {
	save := getRegInt(L)
	setRegInt(L, it)
	err := L.Call(nargs, nresult)
	setRegInt(L, save)
	return err
}

var orgArgHook func(*interpreter.Interpreter, []string) ([]string, error)

var luaArgsFilter lua.Pushable = lua.TNil{}

func newArgHook(it *interpreter.Interpreter, args []string) ([]string, error) {
	L, err := NewNyagosLua()
	if err != nil {
		return nil, err
	}
	defer L.Close()
	L.Push(luaArgsFilter)
	if !L.IsFunction(-1) {
		return orgArgHook(it, args)
	}
	L.NewTable()
	for i := 0; i < len(args); i++ {
		L.PushString(args[i])
		L.RawSetI(-2, lua.Integer(i))
	}
	if err := L.Call(1, 1); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return orgArgHook(it, args)
	}
	if L.GetType(-1) != lua.LUA_TTABLE {
		return orgArgHook(it, args)
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
	return orgArgHook(it, newargs)
}

var orgOnCommandNotFound func(*interpreter.Interpreter, error) error

var luaOnCommandNotFound lua.Pushable = lua.TNil{}

func on_command_not_found(inte *interpreter.Interpreter, err error) error {
	L, err := NewNyagosLua()
	if err != nil {
		return err
	}
	defer L.Close()

	L.Push(luaOnCommandNotFound)
	if !L.IsFunction(-1) {
		L.Pop(1)
		return orgOnCommandNotFound(inte, err)
	}
	L.NewTable()
	for key, val := range inte.Args {
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

var option_table_member = map[string]IProperty{
	"glob": &lua.BoolProperty{&interpreter.WildCardExpansionAlways},
}

func getOption(L lua.Lua) int {
	key, key_err := L.ToString(2)
	if key_err != nil {
		return L.Push(nil, key_err)
	}
	val, val_ok := option_table_member[key]
	if !val_ok {
		return L.Push(nil)
	}
	return L.Push(val)
}

func setOption(L lua.Lua) int {
	key, key_err := L.ToString(2)
	if key_err != nil {
		return L.Push(nil, key_err)
	}
	opt, opt_ok := option_table_member[key]
	if !opt_ok {
		print(key, " not found\n")
		return L.Push(nil)
	}
	if err := opt.Set(L, 3); err != nil {
		return L.Push(nil, err.Error())
	} else {
		return L.Push(true)
	}
}

var nyagos_table_member map[string]lua.Pushable

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
			value, value_err := L.ToPushable(3)
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

var share_table = map[string]lua.Pushable{}

func getShareTable(L lua.Lua) int {
	key, keyErr := L.ToString(-1)
	if keyErr != nil {
		return L.Push(nil, keyErr)
	}
	if value, ok := share_table[key]; ok {
		return L.Push(value)
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
	value, valErr := L.ToPushable(-1)
	if valErr != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", key, valErr.Error())
		return L.Push(nil, valErr)
	}
	share_table[key] = value
	return 1
}

var hook_setuped = false

func NewNyagosLua() (lua.Lua, error) {
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

	if !hook_setuped {
		orgArgHook = interpreter.SetArgsHook(newArgHook)

		orgOnCommandNotFound = interpreter.OnCommandNotFound
		interpreter.OnCommandNotFound = on_command_not_found
		hook_setuped = true
	}
	return this, nil
}

var silentmode = false

func init() {
	nyagos_table_member = map[string]lua.Pushable{
		"access": lua.TGoFunction(cmdAccess),
		"alias": &lua.VirtualTable{
			Name:     "nyagos.alias",
			Index:    cmdGetAlias,
			NewIndex: cmdSetAlias},
		"antihistquot": lua.StringProperty{&history.DisableMarks},
		"argsfilter":   lua.Property{&luaArgsFilter},
		"atou":         lua.TGoFunction(cmdAtoU),
		"key": &lua.VirtualTable{
			Name:     "nyagos.key",
			Index:    cmdGetBindKey,
			NewIndex: cmdBindKey},
		"bindkey":         lua.TGoFunction(cmdBindKey),
		"chdir":           lua.TGoFunction(cmdChdir),
		"commit":          lua.StringProperty{&Commit},
		"commonprefix":    lua.TGoFunction(cmdCommonPrefix),
		"completion_hook": lua.Property{&completionHook},
		"create_object":   lua.TGoFunction(ole.CreateObject),
		"default_prompt":  lua.TGoFunction(nyagosPrompt),
		"elevated":        lua.TGoFunction(cmdElevated),
		"env": &lua.VirtualTable{
			Name:     "nyagos.env",
			Index:    cmdGetEnv,
			NewIndex: cmdSetEnv},
		"eval":         lua.TGoFunction(cmdEval),
		"exec":         lua.TGoFunction(cmdExec),
		"filter":       lua.Property{&luaFilter},
		"getalias":     lua.TGoFunction(cmdGetAlias),
		"getenv":       lua.TGoFunction(cmdGetEnv),
		"gethistory":   lua.TGoFunction(cmdGetHistory),
		"getkey":       lua.TGoFunction(cmdGetKey),
		"getviewwidth": lua.TGoFunction(cmdGetViewWidth),
		"getwd":        lua.TGoFunction(cmdGetwd),
		"glob":         lua.TGoFunction(cmdGlob),
		"goarch":       lua.TString(runtime.GOARCH),
		"goversion":    lua.TString(runtime.Version()),
		"histchar":     lua.StringProperty{&history.Mark},
		"history": &lua.VirtualTable{
			Name:  "nyagos.history",
			Index: cmdGetHistory},
		"lines":                lua.TGoFunction(cmdLines),
		"loadfile":             lua.TGoFunction(cmdLoadFile),
		"netdrivetounc":        lua.TGoFunction(cmdNetDriveToUNC),
		"on_command_not_found": lua.Property{&luaOnCommandNotFound},
		"open":                 lua.TGoFunction(cmdOpenFile),
		"option": &lua.VirtualTable{
			Name:     "nyagos.option",
			Index:    getOption,
			NewIndex: setOption},
		"pathjoin":       lua.TGoFunction(cmdPathJoin),
		"prompt":         lua.Property{&prompt_hook},
		"quotation":      lua.StringProperty{&readline.Delimiters},
		"raweval":        lua.TGoFunction(cmdRawEval),
		"rawexec":        lua.TGoFunction(cmdRawExec),
		"resetcharwidth": lua.TGoFunction(cmdResetCharWidth),
		"setalias":       lua.TGoFunction(cmdSetAlias),
		"setenv":         lua.TGoFunction(cmdSetEnv),
		"setrunewidth":   lua.TGoFunction(cmdSetRuneWidth),
		"shellexecute":   lua.TGoFunction(cmdShellExecute),
		"silentmode":     &lua.BoolProperty{&silentmode},
		"stamp":          lua.StringProperty{&Stamp},
		"stat":           lua.TGoFunction(cmdStat),
		"utoa":           lua.TGoFunction(cmdUtoA),
		"version":        lua.StringProperty{&Version},
		"which":          lua.TGoFunction(cmdWhich),
		"write":          lua.TGoFunction(cmdWrite),
		"writerr":        lua.TGoFunction(cmdWriteErr),
	}
}
