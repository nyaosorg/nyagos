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

extern int luaToGoBridge(lua_State *);
static void gLua_pushbridge(lua_State*L)
{ return lua_pushcfunction(L,luaToGoBridge);}
*/
import "C"

import "errors"
import "fmt"
import "unsafe"

type Lua struct {
	lua *C.lua_State
}

func New() *Lua {
	return &Lua{C.luaL_newstate()}
}

const REGISTORYINDEX = C.LUA_REGISTRYINDEX

//export luaToGoBridge
func luaToGoBridge(lua *C.lua_State) int {
	f, _, _ := lua_touserdata.Call(uintptr(unsafe.Pointer(lua)), 1)
	f_ := *(*goFunctionT)(unsafe.Pointer(f))
	lua_remove.Call(uintptr(unsafe.Pointer(lua)), 1)
	L := Lua{lua}
	return int(f_.function(&L))
}

func (this *Lua) SetField(index int, str string) {
	cstr := C.CString(str)
	C.lua_setfield(this.lua, C.int(index), cstr)
	C.free(unsafe.Pointer(cstr))
}

func (this *Lua) GetField(index int, str string) {
	cstr := C.CString(str)
	C.lua_getfield(this.lua, C.int(index), cstr)
	C.free(unsafe.Pointer(cstr))
}

func (this *Lua) SetGlobal(str string) {
	cstr := C.CString(str)
	C.lua_setglobal(this.lua, cstr)
	C.free(unsafe.Pointer(cstr))
}

func (this *Lua) GetGlobal(str string) {
	cstr := C.CString(str)
	C.lua_getglobal(this.lua, cstr)
	C.free(unsafe.Pointer(cstr))
}

type goFunctionT struct {
	function func(*Lua) int
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

func (this *Lua) ToString(index int) (string, error) {
	if !this.IsString(index) {
		return "", fmt.Errorf("Lua.ToString(%d): Not String", index)
	}
	var length C.size_t
	p := C.lua_tolstring(this.lua, C.int(index), &length)
	if p == nil {
		return "", fmt.Errorf("Lua.ToString(%d): Empty", index)
	} else {
		return C.GoStringN(p, C.int(length)), nil
	}
}

func (this *Lua) ToAnsiString(index int) []byte {
	var length C.size_t
	p := C.lua_tolstring(this.lua, C.int(index), &length)
	if p == nil || length <= 0 {
		return []byte{}
	} else {
		return C.GoBytes(unsafe.Pointer(p), C.int(length))
	}
}

func (this *Lua) Load(fname string) error {
	if C.gLuaL_loadfile(this.lua, C.CString(fname)) != 0 {
		msg, err := this.ToString(-1)
		if err == nil {
			return fmt.Errorf("%s: %s", fname, msg)
		} else {
			return err
		}
	}
	return nil
}

func (this *Lua) Call(nargs, nresult int) error {
	if C.gLua_pcall(this.lua, C.int(nargs), C.int(nresult), 0) != 0 {
		if this.IsString(-1) {
			msg, err := this.ToString(-1)
			if err == nil {
				return errors.New(msg)
			} else {
				return err
			}
		} else {
			return errors.New("<Lua Error>")
		}
	}
	return nil
}

func (this *Lua) NewTable() {
	C.lua_createtable(this.lua, 0, 0)
}

func (this *Lua) PushString(str string) {
	cstr := C.CString(str)
	C.lua_pushstring(this.lua, cstr)
	C.free(unsafe.Pointer(cstr))
}
