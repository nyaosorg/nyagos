package lua

import (
	"context"
	"unsafe"
)

const contextRegKey = "github.com/zetamatta/nyagos/lua"

type contextContainer struct {
	context.Context
}

func (L Lua) context() context.Context {
	L.PushValue(LUA_REGISTRYINDEX)
	typ := L.GetField(-1, contextRegKey)
	defer L.Pop(2)
	if typ != LUA_TUSERDATA {
		return nil
	}
	cc := (*contextContainer)(L.ToUserData(-1))
	if cc != nil {
		return cc.Context
	} else {
		return nil
	}
}

func (L Lua) Context() context.Context {
	ctx := L.context()
	if ctx == nil {
		println("(lua.Lua)Context(): could not found context object")
	}
	return ctx
}

func (L Lua) SetContext(ctx context.Context) {
	L.PushValue(LUA_REGISTRYINDEX)
	if ctx != nil {
		c := contextContainer{Context: ctx}
		L.NewUserDataFrom(unsafe.Pointer(&c), unsafe.Sizeof(c))
	} else {
		L.PushNil()
	}
	L.SetField(-2, contextRegKey)
	L.Pop(1)
}

func (L Lua) CallWithContext(ctx context.Context, nargs, nresult int) error {
	save := L.context()
	defer L.SetContext(save)

	L.SetContext(ctx)
	return L.Call(nargs, nresult)
}
