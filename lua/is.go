package lua

var lua_type = luaDLL.NewProc("lua_type")

func (this Lua) GetType(index int) int {
	rv, _, _ := lua_type.Call(this.State(), uintptr(index))
	return int(rv)
}

// 'lua_isfunction' is implemented as C-macro in the header file.
func (this Lua) IsFunction(index int) bool {
	return this.GetType(index) == LUA_TFUNCTION
}

// 'lua_istable' is implemented as C-macro in the header file.
func (this Lua) IsTable(index int) bool {
	return this.GetType(index) == LUA_TTABLE
}

var lua_isstring = luaDLL.NewProc("lua_isstring")

func (this Lua) IsString(index int) bool {
	rc, _, _ := lua_isstring.Call(this.State(), uintptr(index))
	return rc != 0
}
