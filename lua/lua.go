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
#include "lua.h"
#include "lualib.h"
#include "lauxlib.h"

<<<<<<< HEAD
static int gLua_pcall(lua_State* L,int x,int y,int z)
{ return lua_pcall(L,x,y,z); }

static const char *gLua_tostring(lua_State* L,int i)
{ return lua_tostring(L,i); }

static int gLuaL_loadfile(lua_State* L,const char *filename)
{ return luaL_loadfile(L,filename); }

static void gLua_pushcfunction(lua_State* L,lua_CFunction f)
{
	lua_pushcfunction(L,f);
}

extern int LuaAlias(lua_State*);
extern int LuaSetEnv(lua_State*);

static int setfunctions(lua_State* L)
{
	lua_pushcfunction(L,LuaAlias);
	lua_setglobal(L,"alias");
	lua_pushcfunction(L,LuaSetEnv);
	lua_setglobal(L,"setenv");
}

*/
import "C"
import "errors"
import "strings"
import "os"

import "../alias/table"

type Lua struct {
	lua *C.lua_State
}

func NewLua() *Lua {
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

func (this *Lua) ToString(i int) string {
	return C.GoString(C.gLua_tostring(this.lua, C.int(i)))
}

func (this *Lua) Load(fname string) error {
	if C.gLuaL_loadfile(this.lua, C.CString(fname)) != 0 {
		return errors.New(fname + ": " + this.ToString(-1))
	}
	return nil
}

func (this *Lua) Call(fname string) error {
	if err := this.Load(fname); err != nil {
		return err
	}
	if C.gLua_pcall(this.lua, 0, 0, 0) != 0 {
		return errors.New(fname + ": " + this.ToString(-1))
	}
	return nil
}
