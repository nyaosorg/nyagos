package lua

import (
	"errors"
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

func (this Lua) testUData(index int, tname string) (uintptr, error) {
	// print("TestUData(", index, ",'", tname, "')\n")
	tname_ptr, tname_err := syscall.BytePtrFromString(tname)
	if tname_err != nil {
		return 0, tname_err
	}
	rv, _, _ := luaL_testudata.Call(this.State(), uintptr(index), uintptr(unsafe.Pointer(tname_ptr)))
	return rv, nil
}

var ErrTestUData = errors.New("Failed to TestUData")
var noOperation = func() {}

func (this Lua) TestUDataTo(index int, tname string, p interface{}) (func(), error) {
	src, err := this.testUData(index, tname)
	if err != nil {
		return noOperation, err
	}
	if src == 0 {
		return noOperation, ErrTestUData
	}
	return this.ToUserDataTo(index, p), nil
}
