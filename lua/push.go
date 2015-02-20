package lua

import (
	"syscall"
	"unsafe"
)

var lua_pushnil = luaDLL.NewProc("lua_pushnil")

func (this *Lua) PushNil() {
	lua_pushnil.Call(this.State())
}

var lua_pushboolean = luaDLL.NewProc("lua_pushboolean")

func (this *Lua) PushBool(value bool) {
	if value {
		lua_pushboolean.Call(this.State(), 1)
	} else {
		lua_pushboolean.Call(this.State(), 0)
	}
}

var lua_pushinteger = luaDLL.NewProc("lua_pushinteger")

func (this *Lua) PushInteger(value Integer) {
	params := make([]uintptr, 0, 4)
	params = append(params, this.State())
	params = value.Expand(params)
	lua_pushinteger.Call(params...)
}

var lua_pushlstring = luaDLL.NewProc("lua_pushlstring")

func (this *Lua) PushAnsiString(data []byte) {
	if data != nil && len(data) > 0 {
		lua_pushlstring.Call(this.State(),
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(len(data)))
	} else {
		this.PushString("")
	}
}

var lua_pushstring = luaDLL.NewProc("lua_pushstring")

func (this *Lua) PushString(str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_pushstring.Call(this.State(), uintptr(unsafe.Pointer(cstr)))
}

var lua_pushlightuserdata = luaDLL.NewProc("lua_pushlightuserdata")

func (this *Lua) PushLightUserData(p unsafe.Pointer) {
	lua_pushlightuserdata.Call(this.State(), uintptr(p))
}

var lua_pushvalue = luaDLL.NewProc("lua_pushvalue")

func (this *Lua) PushValue(index int) {
	lua_pushvalue.Call(this.State(), uintptr(index))
}

func luaToGoBridge(lua uintptr) int {
	f, _, _ := lua_touserdata.Call(lua, 1)
	f_ := *(*goFunctionT)(unsafe.Pointer(f))
	lua_remove_Call(lua, 1)
	L := Lua{lua}
	return int(f_.function(&L))
}

type goFunctionT struct {
	function func(*Lua) int
}

var lua_pushcclosure = luaDLL.NewProc("lua_pushcclosure")

func (this *Lua) PushGoFunction(f func(L *Lua) int) {
	f_ := goFunctionT{f}
	voidptr := this.NewUserData(unsafe.Sizeof(f_))
	*(*goFunctionT)(voidptr) = f_
	this.NewTable()
	lua_pushcclosure.Call(this.State(),
		syscall.NewCallbackCDecl(luaToGoBridge),
		0)
	this.SetField(-2, "__call")
	this.SetMetaTable(-2)
}

func (this *Lua) Push(value interface{}) {
	switch t := value.(type) {
	case nil:
		this.PushNil()
	case bool:
		this.PushBool(t)
	case int:
		this.PushInteger(Integer(t))
	case string:
		this.PushString(t)
	case func(L *Lua) int:
		this.PushGoFunction(t)
	case []byte:
		this.PushAnsiString(t)
	case map[string]interface{}:
		this.NewTable()
		for key, val := range t {
			this.PushString(key)
			this.Push(val)
			this.SetTable(-3)
		}
	default:
		panic("lua.Lua.Push(value): value is not supported type")
	}
}
