package mains

import (
	"errors"
	"fmt"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/yuin/gopher-lua"
)

var initializedRequired = true

func lerror(L Lua, s string) int {
	println(s)
	L.Push(lua.LNil)
	L.Push(lua.LString(s))
	return 2
}

type capsuleT struct {
	Data *ole.IDispatch
}

type methodT struct {
	Name string
	Data *ole.IDispatch
}

func (c capsuleT) ToLValue(L Lua) lua.LValue {
	ud := L.NewUserData()
	ud.Value = &c
	meta := L.NewTable()
	L.SetField(meta, "__gc", L.NewFunction(gc))
	L.SetField(meta, "__index", L.NewFunction(index))
	L.SetField(meta, "__newindex", L.NewFunction(set))
	L.SetMetatable(ud, meta)
	return ud
}

func gc(L Lua) int {
	ud := L.ToUserData(1)
	p, ok := ud.Value.(*capsuleT)
	if !ok {
		return lerror(L, "gc: not capsult_t instance")
	}
	if p.Data != nil {
		p.Data.Release()
		p.Data = nil
	}
	return 0
}

func lua2interface(L Lua, index int) (interface{}, error) {
	valueTmp := L.Get(index)
	if valueTmp == lua.LNil {
		return nil, nil
	} else if valueTmp == lua.LTrue {
		return true, nil
	} else if valueTmp == lua.LFalse {
		return false, nil
	}
	switch value := valueTmp.(type) {
	default:
		return nil, errors.New("lua2interface: not support type")
	case lua.LString:
		return string(value), nil
	case lua.LNumber:
		return int(value), nil
	case *lua.LUserData:
		c, ok := value.Value.(*capsuleT)
		if !ok {
			return nil, errors.New("lua2interface: not a OBJECT")
		}
		return c.Data, nil
	}
}

func lua2interfaceS(L Lua, start, end int) ([]interface{}, error) {
	result := make([]interface{}, end-start+1)
	for i := start; i <= end; i++ {
		val, err := lua2interface(L, i)
		if err != nil {
			return nil, err
		}
		result[i-start] = val
	}
	return result, nil
}

// this:_call("METHODNAME",params...)
func call1(L Lua) int {
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok { // OBJECT_T
		return lerror(L, "call1: not found object")
	}
	p, ok := ud.Value.(*capsuleT)
	if !ok {
		return lerror(L, "call1: not found capsuleT")
	}
	name, ok := L.Get(2).(lua.LString)
	if !ok {
		return lerror(L, "call1: not found methodname")
	}
	return callCommon(L, p.Data, string(name))
}

// this:METHODNAME(params...)
func call2(L Lua) int {
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok {
		return lerror(L, "call2: not found userdata for methodT")
	}
	method, ok := ud.Value.(*methodT)
	if !ok || method.Name == "" {
		return lerror(L, "call2: not found methodT")
	}
	ud, ok = L.Get(2).(*lua.LUserData)
	if !ok {
		return lerror(L, "call2: not found userdata for object_t")
	}
	obj, ok := ud.Value.(*capsuleT)
	if !ok {
		if method.Data == nil {
			return lerror(L, "call2: receiver is not found")
		}
		return callCommon(L, method.Data, method.Name)
		// this code enables `OLEOBJ.PROPERTY.PROPERTY:METHOD()`
	}
	if obj.Data == nil {
		return lerror(L, "call2: OLEOBJECT(): the receiver is null")
	}
	return callCommon(L, obj.Data, method.Name)
}

func callCommon(L Lua, com1 *ole.IDispatch, name string) int {
	count := L.GetTop()
	params, err := lua2interfaceS(L, 3, count)
	if err != nil {
		return lerror(L, fmt.Sprintf("callCommon: %s", err.Error()))
	}
	result, err := oleutil.CallMethod(com1, name, params...)
	if err != nil {
		return lerror(L, fmt.Sprintf("oleutil.CallMethod(%s): %s", name, err.Error()))
	}
	if result.VT == ole.VT_DISPATCH {
		L.Push(capsuleT{result.ToIDispatch()}.ToLValue(L))
	} else {
		L.Push(interfaceToLValue(L, result.Value()))
	}
	return 1
}

func set(L Lua) int {
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok {
		return lerror(L, "set: the 1st argument is not usedata")
	}
	p, ok := ud.Value.(*capsuleT)
	if !ok {
		return lerror(L, "set: the 1st argument is not *capsuleT")
	}
	name, ok := L.Get(2).(lua.LString)
	if !ok {
		return lerror(L, "set: the 2nd argument is not string")
	}
	key, err := lua2interfaceS(L, 3, L.GetTop())
	if err != nil {
		return lerror(L, fmt.Sprintf("set: %s", err.Error()))
	}
	oleutil.PutProperty(p.Data, string(name), key...)
	L.Push(lua.LTrue)
	L.Push(lua.LNil)
	return 2
}

func get(L Lua) int {
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok {
		return lerror(L, "get: 1st argument is not a userdata.")
	}
	p, ok := ud.Value.(*capsuleT)
	if !ok {
		return lerror(L, "get: 1st argument is not *capsuleT")
	}

	name, ok := L.Get(2).(lua.LString)
	if !ok {
		return lerror(L, "get: 2nd argument is not string")
	}

	key, err := lua2interfaceS(L, 3, L.GetTop())
	if err != nil {
		return lerror(L, fmt.Sprintf("get: %s", err.Error()))
	}
	result, err := oleutil.GetProperty(p.Data, string(name), key...)
	if err != nil {
		return lerror(L, fmt.Sprintf("oleutil.GetProperty: %s", err.Error()))
	}
	if result.VT == ole.VT_DISPATCH {
		L.Push(capsuleT{result.ToIDispatch()}.ToLValue(L))
	} else {
		L.Push(interfaceToLValue(L, result.Value()))
	}
	return 1
}

func indexSub(L Lua, thisIndex int, nameIndex int) int {
	name, ok := L.Get(nameIndex).(lua.LString)
	if !ok {
		return lerror(L, "indexSub: not a string")
	}
	switch string(name) {
	case "_call":
		L.Push(L.NewFunction(call1))
		L.Push(lua.LNil)
		return 2
	case "_set":
		L.Push(L.NewFunction(set))
		L.Push(lua.LNil)
		return 2
	case "_get":
		L.Push(L.NewFunction(get))
		L.Push(lua.LNil)
		return 2
	default:
		m := &methodT{Name: string(name)}
		if ud, ok := L.Get(thisIndex).(*lua.LUserData); ok {
			if p, ok := ud.Value.(*capsuleT); ok {
				m.Data = p.Data
			}
		}
		ud := L.NewUserData()
		ud.Value = m

		meta := L.NewTable()
		L.SetField(meta, "__newidex", L.NewFunction(set))
		L.SetField(meta, "__call", L.NewFunction(call2))
		L.SetField(meta, "__index", L.NewFunction(get2))
		L.SetMetatable(ud, meta)
		L.Push(ud)

		return 1
	}
}

func index(L Lua) int {
	return indexSub(L, 1, 2)
}

// THIS.member.member
func get2(L Lua) int {
	ud, ok := L.Get(1).(*lua.LUserData)
	if !ok {
		return lerror(L, "get2: not a userdata")
	}
	m, ok := ud.Value.(*methodT)
	if !ok {
		return lerror(L, "get: not a methodT")
	}
	result, err := oleutil.GetProperty(m.Data, m.Name)
	if err != nil {
		return lerror(L, fmt.Sprintf("oleutil.GetProperty: %s", err.Error()))
	}
	if result.VT == ole.VT_DISPATCH {
		L.Push(capsuleT{result.ToIDispatch()}.ToLValue(L))
	} else {
		L.Push(interfaceToLValue(L, result.Value()))
	}
	return indexSub(L, 3, 2)
}

// CreateObject creates Lua-Object to access COM
func CreateObject(L Lua) int {
	if initializedRequired {
		ole.CoInitialize(0)
		initializedRequired = false
	}
	name, ok := L.Get(1).(lua.LString)
	if !ok {
		return lerror(L, "CreateObject: parameter not a string")
	}
	unknown, err := oleutil.CreateObject(string(name))
	if err != nil {
		return lerror(L, fmt.Sprintf("oleutil.CreateObject: %s", err.Error()))
	}
	obj, err := unknown.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return lerror(L, fmt.Sprintf("unknown.QueryInterfce: %s", err.Error()))
	}
	L.Push(capsuleT{obj}.ToLValue(L))
	return 1
}
