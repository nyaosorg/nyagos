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
	userdata, _, _ := lua_touserdata.Call(lua, 1)
	fn := *(*func(*Lua) int)(unsafe.Pointer(userdata))
	lua_remove_Call(lua, 1)
	return int(fn(&Lua{lua}))
}

var lua_pushcclosure = luaDLL.NewProc("lua_pushcclosure")

func (this *Lua) PushGoFunction(fn func(*Lua) int) {
	userdata := this.NewUserData(unsafe.Sizeof(fn))
	*(*func(*Lua) int)(userdata) = fn
	this.NewTable()
	lua_pushcclosure.Call(this.State(),
		syscall.NewCallbackCDecl(luaToGoBridge),
		0)
	this.SetField(-2, "__call")
	this.SetMetaTable(-2)
}

func (this *Lua) Push(values ...interface{}) int {
	for _, value := range values {
		if value == nil {
			this.PushNil()
			continue
		}
		switch t := value.(type) {
		case bool:
			this.PushBool(t)
		case Integer:
			this.PushInteger(Integer(t))
		case int:
			this.PushInteger(Integer(t))
		case int64:
			this.PushInteger(Integer(t))
		case string:
			this.PushString(t)
		case func(L *Lua) int:
			this.PushGoFunction(t)
		case []byte:
			this.PushAnsiString(t)
		case error:
			this.PushString(t.Error())
		case map[string]interface{}:
			this.NewTable()
			for key, val := range t {
				this.PushString(key)
				this.Push(val)
				this.SetTable(-3)
			}
		case unsafe.Pointer:
			this.PushLightUserData(t)
		default:
			panic("lua.Lua.Push(value): value is not supported type")
		}
	}
	return len(values)
}
