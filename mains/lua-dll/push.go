package lua

import (
	"fmt"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

var lua_pushnil = luaDLL.NewProc("lua_pushnil")

func (this Lua) PushNil() {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushNil()\n", this)
	}
	lua_pushnil.Call(this.State())
}

var lua_pushboolean = luaDLL.NewProc("lua_pushboolean")

func (this Lua) PushBool(value bool) {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushBool(%v)\n", this, value)
	}
	if value {
		lua_pushboolean.Call(this.State(), 1)
	} else {
		lua_pushboolean.Call(this.State(), 0)
	}
}

var lua_pushinteger = luaDLL.NewProc("lua_pushinteger")

func (this Lua) PushInteger(value Integer) {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushInteger(%v)\n", this, value)
	}
	params := make([]uintptr, 0, 4)
	params = append(params, this.State())
	params = value.Expand(params)
	lua_pushinteger.Call(params...)
}

var lua_pushlstring = luaDLL.NewProc("lua_pushlstring")

func (this Lua) PushBytes(data []byte) {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushBytes(len=%v)\n", this, len(data))
	}
	if data != nil && len(data) >= 1 {
		lua_pushlstring.Call(this.State(),
			uintptr(unsafe.Pointer(&data[0])),
			uintptr(len(data)))
	} else {
		zerobyte := []byte{'\000'}
		lua_pushlstring.Call(this.State(),
			uintptr(unsafe.Pointer(&zerobyte[0])),
			0)
	}
}

var lua_pushstring = luaDLL.NewProc("lua_pushstring")

func (this Lua) PushString(str string) {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushString(%v)\n", this, str)
	}
	// BytePtrFromString can not use the string which contains NUL
	array := make([]byte, len(str)+1)
	copy(array, str)
	lua_pushlstring.Call(this.State(),
		uintptr(unsafe.Pointer(&array[0])),
		uintptr(len(str)))
}

var lua_pushlightuserdata = luaDLL.NewProc("lua_pushlightuserdata")

func (this Lua) PushLightUserData(p unsafe.Pointer) {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushLightUserData(%v)", this, p)
	}
	lua_pushlightuserdata.Call(this.State(), uintptr(p))
}

var lua_pushvalue = luaDLL.NewProc("lua_pushvalue")

func (this Lua) PushValue(index int) {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushValue(%v)\n", this, index)
	}
	lua_pushvalue.Call(this.State(), uintptr(index))
}

var lua_pushcclosure = luaDLL.NewProc("lua_pushcclosure")

func (this Lua) PushGoClosure(fn func(Lua) int, n uintptr) {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushGoClosure(%v)\n", this, n)
	}
	lua_pushcclosure.Call(this.State(), syscall.NewCallbackCDecl(fn), n)
}

func (this Lua) PushGoFunction(fn func(Lua) int) {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushGoFunction\n", this)
	}
	this.PushGoClosure(fn, 0)
}

func UpValueIndex(i int) int {
	return LUA_REGISTRYINDEX - i
}

type TGoFunction func(Lua) int

func (this TGoFunction) Push(L Lua) int {
	if trace {
		fmt.Fprintf(os.Stderr, "lua.TGoFunction(%v).Push(%v)\n", this, L)
	}
	L.PushGoFunction(this)
	return 1
}

func (this Lua) PushCFunction(fn uintptr) {
	if trace {
		fmt.Fprintf(os.Stderr, "lua(%v).PushCFunction(%v)\n", this, fn)
	}
	lua_pushcclosure.Call(this.State(), fn, 0)
}

type Object interface {
	Push(Lua) int
}

func (this Lua) Push(values ...interface{}) int {
	if trace {
		fmt.Fprintf(os.Stderr, "lua(%v).Push(%v...)\n", this, values)
	}
	for _, value := range values {
		if value == nil {
			this.PushNil()
			continue
		}
		switch t := value.(type) {
		case bool:
			this.PushBool(t)
		case int:
			this.PushInteger(Integer(t))
		case int64:
			this.PushInteger(Integer(t))
		case string:
			this.PushString(t)
		case func(L Lua) int:
			this.PushGoFunction(t)
		case []byte:
			this.PushBytes(t)
		case error:
			this.PushString(t.Error())
		case TTable:
			this.NewTable()
			for key, val := range t.Dict {
				this.PushString(key)
				this.Push(val)
				this.SetTable(-3)
			}
			for key, val := range t.Array {
				this.Push(key)
				this.Push(val)
				this.SetTable(-3)
			}
		case unsafe.Pointer:
			this.PushLightUserData(t)
		case Object:
			t.Push(this)
		default:
			if !this.PushReflect(value) {
				panic(fmt.Sprintf(
					"lua.Lua.Push(%T): value is not supported type", t))
			}
		}
	}
	return len(values)
}

func (this Lua) PushReflect(value interface{}) bool {
	if trace {
		fmt.Fprintf(os.Stderr, "Lua(%v).PushReflect(%v)\n", this, value)
	}
	if value == nil {
		this.PushNil()
	}
	return this.pushReflect(reflect.ValueOf(value))
}

func (this Lua) pushReflect(value reflect.Value) bool {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		this.PushInteger(Integer(value.Int()))
	case reflect.Uint, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr:
		this.PushInteger(Integer(value.Uint()))
	case reflect.Bool:
		this.PushBool(value.Bool())
	case reflect.String:
		this.PushString(value.String())
	case reflect.Interface:
		this.Push(value.Interface())
	case reflect.Slice, reflect.Array:
		elem := value.Type().Elem()
		if elem.Kind() == reflect.Uint8 {
			buffer := make([]byte, 0, value.Len())
			for i, end := 0, value.Len(); i < end; i++ {
				buffer = append(buffer, byte(value.Index(i).Uint()))
			}
			this.PushBytes(buffer)
		} else {
			this.NewTable()
			for i, end := 0, value.Len(); i < end; i++ {
				val := value.Index(i)
				this.PushInteger(Integer(i + 1))
				this.pushReflect(val)
				this.SetTable(-3)
			}
		}
	case reflect.Map:
		this.NewTable()
		for _, key := range value.MapKeys() {
			this.pushReflect(key)
			val := value.MapIndex(key)
			this.pushReflect(val)
			this.SetTable(-3)
		}
	default:
		return false
	}
	return true
}
