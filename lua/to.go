package lua

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

const dbg = false

var ClosureIsNotAvaliable = errors.New("Can't assign a closure")

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

func (this Lua) ToBytes(index int) []byte {
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

type TString string

func (this TString) Push(L Lua) int {
	L.PushString(string(this))
	return 1
}

var lua_touserdata = luaDLL.NewProc("lua_touserdata")

func (this Lua) ToUserData(index int) unsafe.Pointer {
	rv, _, _ := lua_touserdata.Call(this.State(), uintptr(index))
	return unsafe.Pointer(rv)
}

func PtrAndSize(p interface{}) (uintptr, uintptr) {
	value := reflect.ValueOf(p)
	size := value.Type().Elem().Size()
	return value.Pointer(), size
}

func (this Lua) ToUserDataTo(index int, p interface{}) func() {
	src, _, _ := lua_touserdata.Call(this.State(), uintptr(index))
	dst, size := PtrAndSize(p)
	copyMemory(dst, src, size)
	return func() {
		copyMemory(src, dst, size)
	}
}

var lua_toboolean = luaDLL.NewProc("lua_toboolean")

func (this Lua) ToBool(index int) bool {
	rv, _, _ := lua_toboolean.Call(this.State(), uintptr(index))
	return rv != 0
}

type TRawString []byte

func (this TRawString) String() (string, error) {
	if len(this) <= 0 {
		return "", nil
	} else {
		return string(this), nil
	}
}

func (this TRawString) Push(L Lua) int {
	L.PushBytes(this)
	return 1
}

var lua_tocfunction = luaDLL.NewProc("lua_tocfunction")

func (this Lua) ToCFunction(index int) uintptr {
	rc, _, _ := lua_tocfunction.Call(this.State(), uintptr(index))
	return rc
}

type TCFunction uintptr

func (this TCFunction) Push(L Lua) int {
	L.PushCFunction(uintptr(this))
	return 1
}

type TLuaFunction []byte

func (this TLuaFunction) Push(L Lua) int {
	if L.LoadBufferX("(annonymous)", this, "b") != nil {
		return 1
	} else {
		return 0
	}
}

type TLightUserData struct {
	Data unsafe.Pointer
}

func (this TLightUserData) Push(L Lua) int {
	L.PushLightUserData(this.Data)
	return 1
}

type TFullUserData []byte

func (this TFullUserData) Push(L Lua) int {
	L.NewUserDataFrom(unsafe.Pointer(&this[0]), uintptr(len([]byte(this))))
	return 1
}

var lua_next = luaDLL.NewProc("lua_next")

func (this Lua) Next(index int) int {
	rc, _, _ := lua_next.Call(this.State(), uintptr(index))
	return int(rc)
}

var lua_rawlen = luaDLL.NewProc("lua_rawlen")

func (this Lua) RawLen(index int) uintptr {
	size, _, _ := lua_rawlen.Call(this.State(), uintptr(index))
	return size
}

type MetaTableOwner struct {
	Body Object
	Meta *TTable
}

func (this *MetaTableOwner) Push(L Lua) int {
	this.Body.Push(L)
	if nameObj, nameObj_ok := this.Meta.Dict["__name"]; nameObj_ok {
		if name, name_ok := nameObj.(TRawString); name_ok {
			if dbg {
				print("found meta-name: ", string(name), "\n")
			}
			L.NewMetaTable(string(name))
			this.Meta.PushWithoutNewTable(L)
		} else {
			if dbg {
				print("found meta-name, but could not cast\n")
			}
			this.Meta.Push(L)
		}
	} else {
		if dbg {
			print("not meta table\n")
		}
		this.Meta.Push(L)
	}
	L.SetMetaTable(-2)
	return 1
}

type TTable struct {
	Dict  map[string]Object
	Array map[int]Object
}

func (this TTable) PushWithoutNewTable(L Lua) int {
	for key, val := range this.Dict {
		L.PushString(key)
		val.Push(L)
		L.SetTable(-3)
	}
	for key, val := range this.Array {
		L.Push(key)
		val.Push(L)
		L.SetTable(-3)
	}
	return 1
}

func (this TTable) Push(L Lua) int {
	L.NewTable()
	return this.PushWithoutNewTable(L)
}

func (this Lua) ForInDo(index int, proc func(Lua) error) error {
	this.PushNil() // set first key as nil
	if index < 0 {
		index--
	}
	for this.Next(index) != 0 {
		// Next push KEY and VAL
		err := proc(this)
		/* removes 'value'; keeps 'key' for next iteration */
		this.Pop(1)
		if err != nil {
			this.Pop(1)
			return err
		}
		/*
			While traversing a table, do not call lua_tolstring
			directly on a key, unless you know that the key is
			actually a string. Recall that lua_tolstring may
			change the value at the given index; this confuses
			the next call to lua_next.
		*/
	}
	return nil
}

func (this Lua) ToTable(index int) (*TTable, error) {
	table := make(map[string]Object)
	array := make(map[int]Object)

	err := this.ForInDo(index, func(this Lua) error {
		key, err := this.ToObject(-2)
		if err != nil {
			return err
		}
		val, err := this.ToObject(-1)
		if err != nil {
			return err
		}
		switch t := key.(type) {
		case TString:
			table[string(t)] = val
		case TRawString:
			table[string(t)] = val
		case Integer:
			array[int(t)] = val
		case nil:
			table[""] = val
		}
		return nil
	})
	return &TTable{Dict: table, Array: array}, err
}

type TBool struct {
	Value bool
}

func (this TBool) Push(L Lua) int {
	L.PushBool(this.Value)
	return 1
}

type TNil struct{}

func (this TNil) Push(L Lua) int {
	L.PushNil()
	return 1
}

var NG_UPVALUE_NAME = map[string]struct{}{}

func (this Lua) ToObject(index int) (Object, error) {
	seek_metatable := false
	var err error = nil
	var result Object
	switch this.GetType(index) {
	case LUA_TBOOLEAN:
		result = TBool{this.ToBool(index)}
	case LUA_TFUNCTION:
		if p := this.ToCFunction(index); p != 0 {
			// CFunction
			result = TCFunction(p)
		} else {
			// LuaFunction
			upvalues := this.GetUpValues(index)
			for _, u := range upvalues {
				if _, ok := NG_UPVALUE_NAME[u.Name]; ok {
					if dbg {
						print(u.Name, ":", this.TypeName(u.Type), "\n")
					}
					return nil, ClosureIsNotAvaliable
				}
			}
			this.PushValue(index)
			result = TLuaFunction(this.Dump())
			this.Pop(1)
		}
	case LUA_TLIGHTUSERDATA:
		result = TLightUserData{Data: this.ToUserData(index)}
		seek_metatable = true
	case LUA_TNIL:
		result = TNil{}
	case LUA_TNUMBER:
		var int_result int
		int_result, err = this.ToInteger(index)
		result = Integer(int_result)
	case LUA_TSTRING:
		result = TRawString(this.ToBytes(index))
	case LUA_TTABLE:
		result, err = this.ToTable(index)
		seek_metatable = true
	case LUA_TUSERDATA:
		size := this.RawLen(index)
		ptr := this.ToUserData(index)
		result = TFullUserData(CGoBytes(uintptr(ptr), uintptr(size)))
		seek_metatable = true
	default:
		return nil, errors.New("lua.ToSomeThing: Not supported type found.")
	}
	if err != nil {
		return nil, err
	}
	if seek_metatable && this.GetMetaTable(index) {
		metatable, err := this.ToTable(-1)
		defer this.Pop(1)
		if err != nil {
			return nil, err
		}
		result = &MetaTableOwner{Body: result, Meta: metatable}
	}
	return result, nil
}

func (this Lua) ToInterface(index int) (interface{}, error) {
	t := this.GetType(index)
	switch t {
	case LUA_TBOOLEAN:
		return this.ToBool(index), nil
	case LUA_TNIL:
		return nil, nil
	case LUA_TSTRING:
		return this.ToString(index)
	case LUA_TNUMBER:
		intValue, err := this.ToInteger(index)
		if err != nil {
			return nil, err
		}
		return intValue, nil
	case LUA_TTABLE:
		table := map[interface{}]interface{}{}
		err := this.ForInDo(index, func(this Lua) error {
			key, err := this.ToInterface(-2)
			if err != nil {
				return err
			}
			val, err := this.ToInterface(-1)
			if err != nil {
				return err
			}
			table[key] = val
			return nil
		})
		return table, err
	default:
		return nil, fmt.Errorf("Not support Lua type %v", t)
	}
}
