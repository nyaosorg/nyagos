package lua

import (
	"errors"
	"unsafe"
)

var lua_tointegerx = luaDLL.NewProc("lua_tointegerx")

func (this Lua) ToInteger(index int) (int, error) {
	var issucceeded uintptr
	value, _, _ := lua_tointegerx.Call(this.State(), uintptr(index),
		uintptr(unsafe.Pointer(&issucceeded)))
	if issucceeded != 0 {
		return int(value), nil
	} else {
		return 0, errors.New("ToInteger: the value in not integer on the stack")
	}
}

var lua_tolstring = luaDLL.NewProc("lua_tolstring")

func (this Lua) ToAnsiString(index int) []byte {
	var length uintptr
	p, _, _ := lua_tolstring.Call(this.State(),
		uintptr(index),
		uintptr(unsafe.Pointer(&length)))
	if length <= 0 {
		return []byte{}
	} else {
		return CGoBytes(p, length)
	}
}

func (this Lua) ToString(index int) (string, error) {
	var length uintptr
	p, _, _ := lua_tolstring.Call(this.State(),
		uintptr(index),
		uintptr(unsafe.Pointer(&length)))
	return CGoStringN(p, length), nil
}

var lua_touserdata = luaDLL.NewProc("lua_touserdata")

func (this Lua) ToUserData(index int) unsafe.Pointer {
	rv, _, _ := lua_touserdata.Call(this.State(), uintptr(index))
	return unsafe.Pointer(rv)
}

var lua_toboolean = luaDLL.NewProc("lua_toboolean")

func (this Lua) ToBool(index int) bool {
	rv, _, _ := lua_toboolean.Call(this.State(), uintptr(index))
	return rv != 0
}

type TString struct {
	Value []byte
}

func (this *TString) String() (string, error) {
	if len(this.Value) <= 0 {
		return "", nil
	} else {
		return string(this.Value), nil
	}
}

func (this *TString) Push(L Lua) int {
	L.PushAnsiString(this.Value)
	return 1
}

type TFunction struct {
	Chank []byte
}

func (this *TFunction) Push(L Lua) int {
	err := L.LoadBufferX("(anonymous)", this.Chank, "b")
	if err == nil {
		return 1
	} else {
		return 0
	}
}

func (this Lua) ToSomething(index int) (interface{}, error) {
	switch this.GetType(index) {
	case LUA_TBOOLEAN:
		return this.ToBool(index), nil
	case LUA_TFUNCTION:
		return TFunction{Chank: this.Dump()}, nil
	case LUA_TLIGHTUSERDATA:
		return nil, errors.New("lua.ToAnything: LUA_TLIGHTUSERDATA not supported.")
	case LUA_TNIL:
		return nil, nil
	case LUA_TNUMBER:
		return this.ToInteger(index)
	case LUA_TSTRING:
		return TString{this.ToAnsiString(index)}, nil
	case LUA_TTABLE:
		return nil, errors.New("lua.ToAnything: LUA_TTABLE not supported.")
	case LUA_TUSERDATA:
		return nil, errors.New("lua.ToAnything: LUA_TUSERDATA not supported.")
	default:
		return nil, errors.New("lua.ToAnything: Not supported type found.")
	}
}
