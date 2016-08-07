package lua

import (
	"fmt"
	"syscall"
	"unsafe"
)

var lua_pushnil = luaDLL.NewProc("lua_pushnil")

func (this Lua) PushNil() {
	lua_pushnil.Call(this.State())
}

var lua_pushboolean = luaDLL.NewProc("lua_pushboolean")

func (this Lua) PushBool(value bool) {
	if value {
		lua_pushboolean.Call(this.State(), 1)
	} else {
		lua_pushboolean.Call(this.State(), 0)
	}
}

var lua_pushinteger = luaDLL.NewProc("lua_pushinteger")

func (this Lua) PushInteger(value Integer) {
	params := make([]uintptr, 0, 4)
	params = append(params, this.State())
	params = value.Expand(params)
	lua_pushinteger.Call(params...)
}

var lua_pushlstring = luaDLL.NewProc("lua_pushlstring")

func (this Lua) PushBytes(data []byte) {
	if data != nil {
		lua_pushlstring.Call(this.State(),
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(len(data)))
	} else {
		zerobyte := []byte{'\000'}
		lua_pushlstring.Call(this.State(),
			uintptr(unsafe.Pointer(&zerobyte[0])),
			0)
	}
}

var lua_pushstring = luaDLL.NewProc("lua_pushstring")

func (this Lua) PushString(str string) {
	// BytePtrFromString can not use the string which contains NUL
	array := make([]byte, len(str)+1)
	copy(array, str)
	lua_pushlstring.Call(this.State(),
		uintptr(unsafe.Pointer(&array[0])),
		uintptr(len(str)))
}

var lua_pushlightuserdata = luaDLL.NewProc("lua_pushlightuserdata")

func (this Lua) PushLightUserData(p unsafe.Pointer) {
	lua_pushlightuserdata.Call(this.State(), uintptr(p))
}

var lua_pushvalue = luaDLL.NewProc("lua_pushvalue")

func (this Lua) PushValue(index int) {
	lua_pushvalue.Call(this.State(), uintptr(index))
}

var lua_pushcclosure = luaDLL.NewProc("lua_pushcclosure")

func (this Lua) PushGoClosure(fn func(Lua) int, n uintptr) {
	lua_pushcclosure.Call(this.State(), syscall.NewCallbackCDecl(fn), n)
}

func (this Lua) PushGoFunction(fn func(Lua) int) {
	this.PushGoClosure(fn, 0)
}

func UpValueIndex(i int) int {
	return LUA_REGISTRYINDEX - i
}

type TGoFunction func(Lua) int

func (this TGoFunction) Push(L Lua) int {
	L.PushGoFunction(this)
	return 1
}

func (this Lua) PushCFunction(fn uintptr) {
	lua_pushcclosure.Call(this.State(), fn, 0)
}

type Pushable interface {
	Push(Lua) int
}

func (this Lua) Push(values ...interface{}) int {
	for _, value := range values {
		if value == nil {
			this.PushNil()
			continue
		}
		switch t := value.(type) {
		case bool:
			this.PushBool(t)
		case int:
			this.PushInteger(Integer(t))
		case int64:
			this.PushInteger(Integer(t))
		case string:
			this.PushString(t)
		case func(L Lua) int:
			this.PushGoFunction(t)
		case []byte:
			this.PushBytes(t)
		case error:
			this.PushString(t.Error())
		case TTable:
			this.NewTable()
			for key, val := range t.Dict {
				this.PushString(key)
				this.Push(val)
				this.SetTable(-3)
			}
			for key, val := range t.Array {
				this.Push(key)
				this.Push(val)
				this.SetTable(-3)
			}
		case unsafe.Pointer:
			this.PushLightUserData(t)
		case Pushable:
			t.Push(this)
		default:
			panic(fmt.Sprintf("lua.Lua.Push(%T): value is not supported type", t))
		}
	}
	return len(values)
}
