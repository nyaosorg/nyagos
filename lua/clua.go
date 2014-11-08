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

extern int luaToGoBridge(lua_State *);
static void gLua_pushbridge(lua_State*L)
{ return lua_pushcfunction(L,luaToGoBridge);}
*/
import "C"

import "errors"
import "unsafe"

type Lua struct {
	lua *C.lua_State
}

func New() *Lua {
	return &Lua{C.luaL_newstate()}
}

func CGoBytes(p, length uintptr) []byte {
	return C.GoBytes(unsafe.Pointer(p), C.int(length))
}

func CGoStringN(p, length uintptr) string {
	return C.GoStringN((*C.char)(unsafe.Pointer(p)), C.int(length))
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
