package lua

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var lua_load = luaDLL.NewProc("lua_load")

var static_load_buffer [4096]byte

func callback_reader(L uintptr, fd *os.File, size *uintptr) *byte {
	n, err := fd.Read(static_load_buffer[:])
	if err != nil || n == 0 {
		*size = 0
		return nil
	} else {
		*size = uintptr(n)
		return &static_load_buffer[0]
	}
}

func (L Lua) LoadFile(path string, mode string) (int, error) {
	fd, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	path_ptr, path_err := syscall.BytePtrFromString(path)
	if path_err != nil {
		return 0, err
	}
	mode_ptr, mode_err := syscall.BytePtrFromString(mode)
	if mode_err != nil {
		return 0, err
	}
	callback := syscall.NewCallbackCDecl(callback_reader)
	rc, _, _ := lua_load.Call(
		L.State(),
		callback,
		uintptr(unsafe.Pointer(fd)),
		uintptr(unsafe.Pointer(path_ptr)),
		uintptr(unsafe.Pointer(mode_ptr)))

	if rc == LUA_OK {
		return 0, nil
	} else if rc == LUA_ERRSYNTAX {
		return LUA_ERRSYNTAX, errors.New("lua_load: LUA_ERRSYNTAX")
	} else if rc == LUA_ERRMEM {
		return LUA_ERRMEM, errors.New("lua_load: LUA_ERRMEM")
	} else if rc == LUA_ERRGCMM {
		return LUA_ERRGCMM, errors.New("lua_load: LUA_ERRGCMM")
	} else {
		return int(rc), fmt.Errorf("lua_load: returns %d", rc)
	}
}
