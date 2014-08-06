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

static const char *gLua_tostring(lua_State* L,int i)
{ return lua_tostring(L,i); }

static int gLuaL_loadfile(lua_State* L,const char *filename)
{ return luaL_loadfile(L,filename); }

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
	return C.GoString(C.gLua_tostring(this.lua, C.int(index)))
}

func (this *Lua) Load(fname string) error {
	if C.gLuaL_loadfile(this.lua, C.CString(fname)) != 0 {
		return fmt.Errorf("%s: %s", fname, this.ToString(-1))
	}
	return nil
}

func (this *Lua) Call(fname string) error {
	if err := this.Load(fname); err != nil {
		return err
	}
	if C.gLua_pcall(this.lua, 0, 0, 0) != 0 {
		return fmt.Errorf("%s: %s", fname, this.ToString(-1))
	}
	return nil
}

func (this *Lua) NewTable() {
	C.lua_createtable(this.lua, 0, 0)
}

func (this *Lua) PushString(str string) {
	tmp := C.CString(str)
	C.lua_pushstring(this.lua, tmp)
	C.free(unsafe.Pointer(tmp))
}

//export luaToGoBridge
func luaToGoBridge(lua *C.lua_State) int {
	f := *(*goFunctionT)(C.lua_touserdata(lua, 1))
	C.lua_remove(lua, 1)
	L := Lua{lua}
	return int(f.function(&L))
}

func (this *Lua) SetField(index int, str string) {
	tmp := C.CString(str)
	C.lua_setfield(this.lua, C.int(index), tmp)
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
