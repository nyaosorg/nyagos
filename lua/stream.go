package lua

import (
	"syscall"
	"unsafe"

	"../dos/ansicfile"
)

type stream_t struct {
	FilePtr ansicfile.FilePtr
	Closer  uintptr
}

func (L Lua) pushStream(fd ansicfile.FilePtr, closer func(Lua) int) *stream_t {
	var userdata *stream_t
	userdata = (*stream_t)(L.NewUserData(unsafe.Sizeof(*userdata)))
	userdata.FilePtr = fd
	userdata.Closer = syscall.NewCallbackCDecl(closer)
	L.GetField(LUA_REGISTRYINDEX, LUA_FILEHANDLE) // metatable
	L.SetMetaTable(-2)
	return userdata
}

func closer(L Lua) int {
	userdata := (*stream_t)(L.ToUserData(1))
	userdata.FilePtr.Close()
	// print("stream_closed\n")
	return 0
}

func (L Lua) PushStream(filePtr ansicfile.FilePtr) {
	L.pushStream(filePtr, closer)
}

func noncloser(L Lua) int {
	return 0
}

func (L Lua) PushStreamDontClose(filePtr ansicfile.FilePtr) {
	L.pushStream(filePtr, noncloser)
}
