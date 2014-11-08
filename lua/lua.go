package lua

import "errors"
import "fmt"
import "syscall"
import "unsafe"

var luaDLL = syscall.NewLazyDLL("lua52")

func (this *Lua) State() uintptr {
	return uintptr(unsafe.Pointer(this.lua))
}

var luaL_openlibs = luaDLL.NewProc("luaL_openlibs")

func (this *Lua) OpenLibs() {
	luaL_openlibs.Call(this.State())
}

var lua_close = luaDLL.NewProc("lua_close")

func (this *Lua) Close() {
	lua_close.Call(this.State())
}

var lua_isstring = luaDLL.NewProc("lua_isstring")

func (this *Lua) IsString(index int) bool {
	rc, _, _ := lua_isstring.Call(this.State(), uintptr(index))
	return rc != 0
}

func (this *Lua) IsFunction(index int) bool {
	return this.GetType(index) == LUA_TFUNCTION
}

var lua_type = luaDLL.NewProc("lua_type")

func (this *Lua) GetType(index int) int {
	rv, _, _ := lua_type.Call(this.State(), uintptr(index))
	return int(rv)
}

var lua_pushvalue = luaDLL.NewProc("lua_pushvalue")

func (this *Lua) PushValue(index int) {
	lua_pushvalue.Call(this.State(), uintptr(index))
}

func (this *Lua) Source(fname string) error {
	if err := this.Load(fname); err != nil {
		return err
	}
	return this.Call(0, 0)
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

var lua_pushinteger = luaDLL.NewProc("lua_pushinteger")

func (this *Lua) PushInteger(n int) {
	lua_pushinteger.Call(this.State(), uintptr(n))
}

var lua_tointegerx = luaDLL.NewProc("lua_tointegerx")

func (this *Lua) ToInteger(index int) (int, error) {
	var issucceeded uintptr
	value, _, _ := lua_tointegerx.Call(this.State(), uintptr(index),
		uintptr(unsafe.Pointer(&issucceeded)))
	if issucceeded != 0 {
		return int(value), nil
	} else {
		return 0, errors.New("ToInteger: the value in not integer on the stack")
	}
}

var lua_settable = luaDLL.NewProc("lua_settable")

func (this *Lua) SetTable(index int) {
	lua_settable.Call(this.State(), uintptr(index))
}

var lua_gettable = luaDLL.NewProc("lua_gettable")

func (this *Lua) GetTable(index int) {
	lua_gettable.Call(this.State(), uintptr(index))
}

var lua_setmetatable = luaDLL.NewProc("lua_setmetatable")

func (this *Lua) SetMetaTable(index int) {
	lua_setmetatable.Call(this.State(), uintptr(index))
}

var lua_gettop = luaDLL.NewProc("lua_gettop")

func (this *Lua) GetTop() int {
	rv, _, _ := lua_gettop.Call(this.State())
	return int(rv)
}

var lua_settop = luaDLL.NewProc("lua_settop")

func (this *Lua) SetTop(index int) {
	lua_settop.Call(this.State(), uintptr(index))
}

func (this *Lua) Pop(n int) {
	this.SetTop(-n - 1)
}

var lua_pushlightuserdata = luaDLL.NewProc("lua_pushlightuserdata")

func (this *Lua) PushLightUserData(p unsafe.Pointer) {
	lua_pushlightuserdata.Call(this.State(), uintptr(p))
}

var lua_touserdata = luaDLL.NewProc("lua_touserdata")

func (this *Lua) ToUserData(index int) unsafe.Pointer {
	rv, _, _ := lua_touserdata.Call(this.State(), uintptr(index))
	return unsafe.Pointer(rv)
}

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

var lua_toboolean = luaDLL.NewProc("lua_toboolean")

func (this *Lua) ToBool(index int) bool {
	rv, _, _ := lua_toboolean.Call(this.State(), uintptr(index))
	if rv != 0 {
		return true
	} else {
		return false
	}
}

var lua_rawseti = luaDLL.NewProc("lua_rawseti")

func (this *Lua) RawSetI(index int, n int) {
	lua_rawseti.Call(this.State(), uintptr(index), uintptr(n))
}

var lua_remove = luaDLL.NewProc("lua_remove")

func (this *Lua) Remove(index int) {
	lua_remove.Call(this.State(), uintptr(index))
}

var lua_replace = luaDLL.NewProc("lua_replace")

func (this *Lua) Replace(index int) {
	lua_replace.Call(this.State(), uintptr(index))
}

var lua_setglobal = luaDLL.NewProc("lua_setglobal")

func (this *Lua) SetGlobal(str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_setglobal.Call(this.State(), uintptr(unsafe.Pointer(cstr)))
}

var lua_setfield = luaDLL.NewProc("lua_setfield")

func (this *Lua) SetField(index int, str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_setfield.Call(this.State(), uintptr(index), uintptr(unsafe.Pointer(cstr)))
}

var lua_getfield = luaDLL.NewProc("lua_getfield")

func (this *Lua) GetField(index int, str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_getfield.Call(this.State(), uintptr(index), uintptr(unsafe.Pointer(cstr)))
}

var lua_getglobal = luaDLL.NewProc("lua_getglobal")

func (this *Lua) GetGlobal(str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_getglobal.Call(this.State(), uintptr(unsafe.Pointer(cstr)))
}

var lua_createtable = luaDLL.NewProc("lua_createtable")

func (this *Lua) NewTable() {
	lua_createtable.Call(this.State(), 0, 0)
}

var lua_pushstring = luaDLL.NewProc("lua_pushstring")

func (this *Lua) PushString(str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_pushstring.Call(this.State(), uintptr(unsafe.Pointer(cstr)))
}

var lua_tolstring = luaDLL.NewProc("lua_tolstring")

func (this *Lua) ToAnsiString(index int) []byte {
	var length uintptr
	p, _, err := lua_tolstring.Call(this.State(),
		uintptr(index),
		uintptr(unsafe.Pointer(&length)))
	if err != nil || length <= 0 {
		return []byte{}
	} else {
		return CGoBytes(p, length)
	}
}

func (this *Lua) ToString(index int) (string, error) {
	if !this.IsString(index) {
		return "", fmt.Errorf("Lua.ToString(%d): Not String", index)
	}
	var length uintptr
	p, _, _ := lua_tolstring.Call(this.State(),
		uintptr(index),
		uintptr(unsafe.Pointer(&length)))
	return CGoStringN(p, length), nil
}
