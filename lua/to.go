package lua

import (
	"errors"
	"fmt"
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

func (this TFunction) Push(L Lua) int {
	err := L.LoadBufferX("(anonymous)", this.Chank, "b")
	if err == nil {
		return 1
	} else {
		return 0
	}
}

var lua_next = luaDLL.NewProc("lua_next")

func (this Lua) Next(index int) int {
	rc, _, _ := lua_next.Call(this.State(), uintptr(index))
	return int(rc)
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
		top := this.GetTop()
		defer this.SetTop(top)
		table := map[string]interface{}{}
		this.PushNil()
		if index < 0 {
			index--
		}
		for this.Next(index) != 0 {
			key, keyErr := this.ToSomething(-2)
			if keyErr == nil {
				val, valErr := this.ToSomething(-1)
				if valErr == nil {
					switch t := key.(type) {
					case string:
						table[t] = val
					case int:
						table[fmt.Sprintf("%d", t)] = val
					case nil:
						table[""] = val
					}
				}
			}
			this.Pop(1)
		}
		return table, nil
	case LUA_TUSERDATA:
		return nil, errors.New("lua.ToAnything: LUA_TUSERDATA not supported.")
	default:
		return nil, errors.New("lua.ToAnything: Not supported type found.")
	}
}
