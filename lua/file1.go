package lua

import (
	"syscall"
	"unsafe"

	"../dos/ansicfile"
)

type CFileT struct {
	FilePtr ansicfile.FilePtr
	Closer  uintptr
}

func ioClose(L Lua) int {
	userdata := (*CFileT)(L.ToUserData(1))
	userdata.FilePtr.Close()
	return 0
}

func OpenByUtf8Path(L Lua) int {
	path, path_err := L.ToString(-2)
	if path_err != nil {
		return L.Push(nil, path_err.Error())
	}
	mode, mode_err := L.ToString(-1)
	if mode_err != nil {
		return L.Push(nil, mode_err.Error())
	}
	fd, fd_err := ansicfile.Open(path, mode)
	if fd_err != nil {
		return L.Push(nil, fd_err)
	}
	var userdata *CFileT
	userdata = (*CFileT)(L.NewUserData(unsafe.Sizeof(*userdata)))
	userdata.FilePtr = fd
	userdata.Closer = syscall.NewCallbackCDecl(ioClose)

	L.GetField(LUA_REGISTRYINDEX, LUA_FILEHANDLE) // metatable
	L.SetMetaTable(-2)

	L.PushNil()
	return 2
}
