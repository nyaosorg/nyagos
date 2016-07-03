package lua

import (
	"syscall"
	"unsafe"
)

var luaL_newmetatable = luaDLL.NewProc("luaL_newmetatable")

func (this Lua) NewMetaTable(tname string) (uintptr, error) {
	tname_ptr, tname_err := syscall.BytePtrFromString(tname)
	if tname_err != nil {
		return 0, tname_err
	}
	rc, _, _ := luaL_newmetatable.Call(this.State(), uintptr(unsafe.Pointer(tname_ptr)))
	return rc, nil
}

var luaL_testudata = luaDLL.NewProc("luaL_testudata")

func (this Lua) TestUData(index int, tname string) (unsafe.Pointer, error) {
	// print("TestUData(", index, ",'", tname, "')\n")
	tname_ptr, tname_err := syscall.BytePtrFromString(tname)
	if tname_err != nil {
		return nil, tname_err
	}
	rv, _, _ := luaL_testudata.Call(this.State(), uintptr(index), uintptr(unsafe.Pointer(tname_ptr)))
	return unsafe.Pointer(rv), nil
}
