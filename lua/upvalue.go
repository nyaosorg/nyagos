package lua

var lua_typename = luaDLL.NewProc("lua_typename")

func (this Lua) TypeName(tp int) string {
	rc, _, _ := lua_typename.Call(this.State(), uintptr(tp))
	return CGoStringZ(rc)
}

var lua_getupvalue = luaDLL.NewProc("lua_getupvalue")

func (this Lua) GetUpValue_(funcindex, n int) uintptr {
	result, _, _ := lua_getupvalue.Call(this.State(), uintptr(funcindex), uintptr(n))
	return result
}

type UpValue struct {
	Name  string
	Index int
	Value Pushable
}

func (this Lua) GetUpValue(funcindex, n int) (string, bool) {
	pointer := this.GetUpValue_(funcindex, n)
	if pointer == 0 {
		return "", false
	} else {
		return CGoStringZ(pointer), true
	}
}

func (this Lua) GetUpValues(funcindex int) []UpValue {
	values := make([]UpValue, 0)
	for i := 1; ; i++ {
		name, ok := this.GetUpValue(funcindex, i)
		if !ok {
			break
		}
		var value1 Pushable
		switch this.GetType(-1) {
		case LUA_TSTRING:
			str, _ := this.ToString(-1)
			value1 = TString{str}
		case LUA_TNUMBER:
			int_result, _ := this.ToInteger(-1)
			value1 = Integer(int_result)
		case LUA_TFUNCTION:
			if p := this.ToCFunction(-1); p != 0 {
				value1 = TCFunction(p)
			} else {
				value1 = new(TNil)
			}
		}
		values = append(values, UpValue{Name: name, Index: i, Value: value1})
		this.Pop(1)
	}
	return values
}
