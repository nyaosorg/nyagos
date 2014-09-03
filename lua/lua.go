package lua

/*
use
	Win32_dllw4: DLL and import library built with MingW gcc 4.3,
	creates dependency with MSVCRT.DLL
from http://sourceforge.net/projects/luabinaries/files/5.2.3/Windows%20Libraries/Dynamic/
*/

/*
#cgo CFLAGS: -I ./include
#cgo LDFLAGS: -llua52 -L.
#include "stdlib.h"
#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"

static int gLua_pcall(lua_State* L,int x,int y,int z)
{ return lua_pcall(L,x,y,z); }

static int gLuaL_loadfile(lua_State* L,const char *filename)
{ return luaL_loadfile(L,filename); }

static int gLua_isfunction(lua_State* L,int i)
{ return lua_isfunction(L,i); }

static void gLua_pop(lua_State* L,int n)
{ lua_pop(L,n); }

extern int luaToGoBridge(lua_State *);
static void gLua_pushbridge(lua_State*L)
{ return lua_pushcfunction(L,luaToGoBridge);}
*/
import "C"

import "fmt"
import "unsafe"

type Lua struct {
	lua *C.lua_State
}

const Registory = C.LUA_REGISTRYINDEX

type goFunctionT struct {
	function func(*Lua) int
}

func New() *Lua {
	this := new(Lua)
	this.lua = C.luaL_newstate()
	return this
}

func (this *Lua) OpenLibs() {
	C.luaL_openlibs(this.lua)
}

func (this *Lua) Close() {
	C.lua_close(this.lua)
}

func (this *Lua) ToString(index int) string {
	var length C.size_t
	p := C.lua_tolstring(this.lua, C.int(index), &length)
	return C.GoStringN(p, C.int(length))
}

func (this *Lua) ToAnsiString(index int) []byte {
	var length C.size_t
	p := C.lua_tolstring(this.lua, C.int(index), &length)
	return C.GoBytes(unsafe.Pointer(p), C.int(length))
}

func (this *Lua) IsString(index int) bool {
	return C.lua_isstring(this.lua, C.int(index)) != 0
}

func (this *Lua) IsFunction(index int) bool {
	return C.gLua_isfunction(this.lua, C.int(index)) != 0
}

func (this *Lua) PushValue(index int) {
	C.lua_pushvalue(this.lua, C.int(index))
}

func (this *Lua) Load(fname string) error {
	if C.gLuaL_loadfile(this.lua, C.CString(fname)) != 0 {
		return fmt.Errorf("%s: %s", fname, this.ToString(-1))
	}
	return nil
}

func (this *Lua) Call(nargs, nresult int) error {
	if C.gLua_pcall(this.lua, C.int(nargs), C.int(nresult), 0) != 0 {
		return fmt.Errorf("%s", this.ToString(-1))
	}
	return nil
}

func (this *Lua) Source(fname string) error {
	if err := this.Load(fname); err != nil {
		return err
	}
	return this.Call(0, 0)
}

func (this *Lua) NewTable() {
	C.lua_createtable(this.lua, 0, 0)
}

func (this *Lua) PushString(str string) {
	tmp := C.CString(str)
	C.lua_pushstring(this.lua, tmp)
	C.free(unsafe.Pointer(tmp))
}

func (this *Lua) PushAnsiString(data []byte) {
	C.lua_pushlstring(this.lua,
		(*C.char)(unsafe.Pointer(&data[0])),
		C.size_t(len(data)))
}

func (this *Lua) PushInteger(n int) {
	C.lua_pushinteger(this.lua, C.lua_Integer(n))
}

//export luaToGoBridge
func luaToGoBridge(lua *C.lua_State) int {
	f := *(*goFunctionT)(C.lua_touserdata(lua, 1))
	C.lua_remove(lua, 1)
	L := Lua{lua}
	return int(f.function(&L))
}

func (this *Lua) SetTable(index int) {
	C.lua_settable(this.lua, C.int(index))
}

func (this *Lua) SetField(index int, str string) {
	tmp := C.CString(str)
	C.lua_setfield(this.lua, C.int(index), tmp)
	C.free(unsafe.Pointer(tmp))
}

func (this *Lua) GetField(index int, str string) {
	tmp := C.CString(str)
	C.lua_getfield(this.lua, C.int(index), tmp)
	C.free(unsafe.Pointer(tmp))
}

func (this *Lua) SetGlobal(str string) {
	tmp := C.CString(str)
	C.lua_setglobal(this.lua, tmp)
	C.free(unsafe.Pointer(tmp))
}

func (this *Lua) GetGlobal(str string) {
	tmp := C.CString(str)
	C.lua_getglobal(this.lua, tmp)
	C.free(unsafe.Pointer(tmp))
}

func (this *Lua) SetMetaTable(index int) {
	C.lua_setmetatable(this.lua, C.int(index))
}

func (this *Lua) PushGoFunction(f func(L *Lua) int) {
	f_ := goFunctionT{f}
	voidptr := C.lua_newuserdata(this.lua, C.size_t(unsafe.Sizeof(f_)))
	*(*goFunctionT)(voidptr) = f_
	this.NewTable()
	C.gLua_pushbridge(this.lua)
	this.SetField(-2, "__call")
	this.SetMetaTable(-2)
}

func (this *Lua) GetTop() int {
	return int(C.lua_gettop(this.lua))
}

func (this *Lua) SetTop(index int) {
	C.lua_settop(this.lua, C.int(index))
}

func (this *Lua) Pop(n int) {
	C.gLua_pop(this.lua, C.int(n))
}

func (this *Lua) PushLightUserData(p unsafe.Pointer) {
	C.lua_pushlightuserdata(this.lua, p)
}

func (this *Lua) ToUserData(index int) unsafe.Pointer {
	return C.lua_touserdata(this.lua, C.int(index))
}
