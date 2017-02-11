package lua

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

var luaDLL = syscall.NewLazyDLL("lua53.dll")

type Integer int64

const LUAINT_PER_UINTPTR = unsafe.Sizeof(Integer(0)) / unsafe.Sizeof(uintptr(0))

func (value Integer) Expand(list []uintptr) []uintptr {
	for i := uintptr(0); i < LUAINT_PER_UINTPTR; i++ {
		list = append(list, uintptr(value))
		value >>= (8 * unsafe.Sizeof(uintptr(1)))
	}
	return list
}

func (value Integer) Push(L Lua) int {
	L.PushInteger(value)
	return 1
}

type Lua uintptr

var luaL_newstate = luaDLL.NewProc("luaL_newstate")

func New() (Lua, error) {
	if err := luaL_newstate.Find(); err != nil {
		return Lua(0), err
	}
	lua, _, _ := luaL_newstate.Call()
	return Lua(lua), nil
}

func (this Lua) State() uintptr {
	return uintptr(this)
}

var luaL_openlibs = luaDLL.NewProc("luaL_openlibs")

func (this Lua) OpenLibs() {
	luaL_openlibs.Call(this.State())
}

var lua_close = luaDLL.NewProc("lua_close")

func (this Lua) Close() error {
	lua_close.Call(this.State())
	return nil
}

func (this Lua) Source(fname string) error {
	if _, err := this.LoadFile(fname, "bt"); err != nil {
		return err
	}
	return this.Call(0, 0)
}

var lua_settable = luaDLL.NewProc("lua_settable")

func (this Lua) SetTable(index int) {
	lua_settable.Call(this.State(), uintptr(index))
}

var lua_gettable = luaDLL.NewProc("lua_gettable")

func (this Lua) GetTable(index int) {
	lua_gettable.Call(this.State(), uintptr(index))
}

var lua_setmetatable = luaDLL.NewProc("lua_setmetatable")

func (this Lua) SetMetaTable(index int) {
	lua_setmetatable.Call(this.State(), uintptr(index))
}

var lua_getmetatable = luaDLL.NewProc("lua_getmetatable")

func (this Lua) GetMetaTable(index int) bool {
	rc, _, _ := lua_getmetatable.Call(this.State(), uintptr(index))
	return rc != 0
	// true: has a metatable / false: has no metatables.
}

var lua_gettop = luaDLL.NewProc("lua_gettop")

func (this Lua) GetTop() int {
	rv, _, _ := lua_gettop.Call(this.State())
	return int(rv)
}

var lua_settop = luaDLL.NewProc("lua_settop")

func (this Lua) SetTop(index int) {
	lua_settop.Call(this.State(), uintptr(index))
}

func (this Lua) Pop(n uint) {
	this.SetTop(-int(n) - 1)
}

var lua_newuserdata = luaDLL.NewProc("lua_newuserdata")

func (this Lua) NewUserData(size uintptr) unsafe.Pointer {
	area, _, _ := lua_newuserdata.Call(this.State(), size)
	return unsafe.Pointer(area)
}

var lua_rawseti = luaDLL.NewProc("lua_rawseti")

func (this Lua) RawSetI(index int, at Integer) {
	params := make([]uintptr, 0, 4)
	params = append(params, this.State(), uintptr(index))
	params = at.Expand(params)
	lua_rawseti.Call(params...)
}

var lua_rawgeti = luaDLL.NewProc("lua_rawgeti")

func (this Lua) RawGetI(index int, at Integer) {
	params := make([]uintptr, 0, 4)
	params = append(params, this.State(), uintptr(index))
	params = at.Expand(params)
	lua_rawgeti.Call(params...)
}

// 5.2
// var lua_remove = luaDLL.NewProc("lua_remove")
// 5.3
var lua_rotate = luaDLL.NewProc("lua_rotate")

func lua_remove_Call(state uintptr, index int) {
	lua_rotate.Call(state, uintptr(index), ^uintptr(0))
	lua_settop.Call(state, ^uintptr(1)) // ^1 == -2
}

func (this Lua) Remove(index int) {
	// 5.2
	// lua_remove.Call(this.State(), uintptr(index))
	// 5.3
	lua_remove_Call(this.State(), index)
}

var lua_replace = luaDLL.NewProc("lua_replace")

func (this Lua) Replace(index int) {
	lua_replace.Call(this.State(), uintptr(index))
}

var lua_setglobal = luaDLL.NewProc("lua_setglobal")

func (this Lua) SetGlobal(str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_setglobal.Call(this.State(), uintptr(unsafe.Pointer(cstr)))
}

var lua_setfield = luaDLL.NewProc("lua_setfield")

func (this Lua) SetField(index int, str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_setfield.Call(this.State(), uintptr(index), uintptr(unsafe.Pointer(cstr)))
}

var lua_getfield = luaDLL.NewProc("lua_getfield")

func (this Lua) GetField(index int, str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_getfield.Call(this.State(), uintptr(index), uintptr(unsafe.Pointer(cstr)))
}

var lua_getglobal = luaDLL.NewProc("lua_getglobal")

func (this Lua) GetGlobal(str string) {
	cstr, err := syscall.BytePtrFromString(str)
	if err != nil {
		panic(err.Error())
	}
	lua_getglobal.Call(this.State(), uintptr(unsafe.Pointer(cstr)))
}

var lua_createtable = luaDLL.NewProc("lua_createtable")

func (this Lua) NewTable() {
	lua_createtable.Call(this.State(), 0, 0)
}

var luaL_loadfilex = luaDLL.NewProc("luaL_loadfilex")

func (this Lua) LoadFileAnsi(fname string) error {
	cfname, err := syscall.BytePtrFromString(fname)
	if err != nil {
		return err
	}
	rc, _, _ := luaL_loadfilex.Call(this.State(),
		uintptr(unsafe.Pointer(cfname)),
		uintptr(0))
	if rc == 0 {
		return nil
	}
	defer this.Pop(1)
	msg, err := this.ToString(-1)
	if err == nil {
		return fmt.Errorf("%s: %s..", fname, msg)
	} else {
		return err
	}
}

var luaL_loadstring = luaDLL.NewProc("luaL_loadstring")

func (this Lua) LoadString(code string) error {
	codePtr, err := syscall.BytePtrFromString(code)
	if err != nil {
		return err
	}
	rc, _, _ := luaL_loadstring.Call(this.State(), uintptr(unsafe.Pointer(codePtr)))
	if rc == 0 {
		return nil
	}
	defer this.Pop(1)
	msg, err := this.ToString(-1)
	if err == nil {
		return errors.New(msg)
	} else {
		return err
	}
}

var lua_pcallk = luaDLL.NewProc("lua_pcallk")

func (this Lua) Call(nargs, nresult int) error {
	rc, _, _ := lua_pcallk.Call(
		this.State(),
		uintptr(nargs),
		uintptr(nresult),
		0,
		0,
		0)
	if rc == 0 {
		return nil
	}
	defer this.Pop(1)
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

var lua_len = luaDLL.NewProc("lua_len")

func (this Lua) Len(index int) {
	lua_len.Call(this.State(), uintptr(index))
}

var lua_newthread = luaDLL.NewProc("lua_newthread")

func (this Lua) NewThread() Lua {
	newthread, _, _ := lua_newthread.Call(this.State())
	return Lua(newthread)
}

func callback_writer(L, p, sz, ud uintptr) uintptr {
	buffer := (*[]byte)(unsafe.Pointer(ud))
	for i := uintptr(0); i < sz; i++ {
		*buffer = append(*buffer, *(*byte)(unsafe.Pointer(p)))
		p++
	}
	return 0
}

var lua_dump = luaDLL.NewProc("lua_dump")

func (this Lua) Dump() []byte {
	buffer := make([]byte, 0, 1024)

	rc, _, _ := lua_dump.Call(
		this.State(),
		syscall.NewCallbackCDecl(callback_writer),
		uintptr(unsafe.Pointer(&buffer)))
	if rc == 0 {
		return buffer
	} else {
		return nil
	}
}

var luaL_loadbufferx = luaDLL.NewProc("luaL_loadbufferx")

func (this Lua) LoadBufferX(title string, chank []byte, mode string) error {
	if this == 0 {
		return errors.New("lua.LoadBufferX: this is null.")
	}
	title_ptr, title_err := syscall.BytePtrFromString(title)
	if title_err != nil {
		return title_err
	}
	mode_ptr, mode_err := syscall.BytePtrFromString(mode)
	if mode_err != nil {
		return mode_err
	}
	if len(chank) <= 0 {
		chank = []byte{';'}
	}

	rc, _, _ := luaL_loadbufferx.Call(
		this.State(),
		uintptr(unsafe.Pointer(&chank[0])),
		uintptr(len(chank)),
		uintptr(unsafe.Pointer(title_ptr)),
		uintptr(unsafe.Pointer(mode_ptr)))
	if rc == 0 {
		return nil
	} else {
		return fmt.Errorf("LUA_ERROR(%d)", rc)
	}
}
