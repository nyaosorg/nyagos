package mains

import (
	"context"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"github.com/zetamatta/go-box"

	"github.com/zetamatta/nyagos/lua"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/texts"
)

type KeyLuaFuncT struct {
	Chank []byte
}

func getBufferForCallBack(L Lua) (*readline.Buffer, int) {
	if L.GetType(1) != lua.LUA_TTABLE {
		return nil, L.Push(nil, "bindKeyExec: call with : not .")
	}
	L.GetField(1, "buffer")
	if L.GetType(-1) != lua.LUA_TLIGHTUSERDATA {
		return nil, L.Push(nil, "bindKey.Call: invalid object")
	}
	buffer := (*readline.Buffer)(L.ToUserData(-1))
	if buffer == nil {
		return nil, L.Push(nil, "bindKey.Call: invalid member")
	}
	L.Pop(1)
	return buffer, 0
}

func callReplace(L Lua) int {
	buffer, stackRc := getBufferForCallBack(L)
	if buffer == nil {
		return stackRc
	}
	pos, pos_err := L.ToInteger(-2)
	if pos_err != nil {
		return L.Push(nil, pos_err.Error())
	}
	str, str_err := L.ToString(-1)
	if str_err != nil {
		return L.Push(nil, str_err.Error())
	}
	pos_zero_base := pos - 1
	if pos_zero_base > buffer.Length {
		return L.Push(nil, fmt.Errorf(":replace: pos=%d: Too big.", pos))
	}
	buffer.ReplaceAndRepaint(pos_zero_base, str)
	return L.Push(true, nil)
}

func callInsert(L Lua) int {
	buffer, stackRc := getBufferForCallBack(L)
	if buffer == nil {
		return stackRc
	}
	text, textErr := L.ToString(2)
	if textErr != nil {
		return L.Push(nil, textErr)
	}
	buffer.InsertAndRepaint(text)
	return L.Push(true)
}

func callKeyFunc(L Lua) int {
	buffer, stackRc := getBufferForCallBack(L)
	if buffer == nil {
		return stackRc
	}
	key, keyErr := L.ToString(2)
	if keyErr != nil {
		return L.Push(nil, keyErr)
	}
	function, funcErr := readline.GetFunc(key)
	if funcErr != nil {
		return L.Push(nil, funcErr)
	}
	ctx := context.Background()
	switch function.Call(ctx, buffer) {
	case readline.ENTER:
		return L.Push(true, true)
	case readline.ABORT:
		return L.Push(true, false)
	default:
		return L.Push(nil)
	}
}

func callLastWord(L Lua) int {
	this, stack_count := getBufferForCallBack(L)
	if this == nil {
		return stack_count
	}
	word, pos := this.CurrentWord()
	return L.Push(word, pos+1)
}

func callFirstWord(L Lua) int {
	this, stack_count := getBufferForCallBack(L)
	if this == nil {
		return stack_count
	}
	word := texts.FirstWord(this.String())
	return L.Push(word, 0)
}

func callBoxListing(L Lua) int {
	// stack +1: readline.Buffer
	// stack +2: table
	// stack +3: index or value
	this, stack_count := getBufferForCallBack(L)
	if this == nil {
		return stack_count
	}
	fmt.Print("\n")
	list := make([]string, 0, 100)
	for i := 1; ; i++ {
		L.Push(i)     // to +3
		L.GetTable(2) //
		str, err := L.ToString(3)
		if err != nil {
			fmt.Fprintln(os.Stderr, "boxprint: "+err.Error())
			break
		}
		if str == "" {
			break
		}
		L.Pop(1)
		list = append(list, str)
	}
	box.Print(nil, list, os.Stdout)
	this.RepaintAll()
	return 0
}

func (this KeyLuaFuncT) String() string {
	return "(lua function)"
}
func (this *KeyLuaFuncT) Call(ctx context.Context, buffer *readline.Buffer) readline.Result {
	L, ok := ctx.Value(lua.PackageId).(Lua)
	if !ok {
		println("(*mains.KeyLuaFuncT)Call: lua instance not found")
		return readline.CONTINUE
	}
	L.LoadBufferX("", this.Chank, "b")
	pos := -1
	var text strings.Builder
	for i, c := range buffer.Buffer {
		if i >= buffer.Length {
			break
		}
		if i == buffer.Cursor {
			pos = text.Len() + 1
		}
		text.WriteRune(c)
	}
	if pos < 0 {
		pos = text.Len() + 1
	}

	L.Push(
		lua.TTable{
			Dict: map[string]lua.Object{
				"pos":         lua.Integer(pos),
				"text":        lua.TString(text.String()),
				"buffer":      lua.TLightUserData{Data: unsafe.Pointer(buffer)},
				"call":        lua.TGoFunction(callKeyFunc),
				"insert":      lua.TGoFunction(callInsert),
				"replacefrom": lua.TGoFunction(callReplace),
				"lastword":    lua.TGoFunction(callLastWord),
				"firstword":   lua.TGoFunction(callFirstWord),
				"boxprint":    lua.TGoFunction(callBoxListing),
			},
			Array: map[int]lua.Object{},
		})
	if err := L.CallWithContext(ctx, 1, 1); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	switch L.GetType(-1) {
	case lua.LUA_TSTRING:
		str, strErr := L.ToString(-1)
		if strErr == nil {
			buffer.InsertAndRepaint(str)
		}
	case lua.LUA_TBOOLEAN:
		if !L.ToBool(-1) {
			buffer.Buffer = []rune{}
			buffer.Length = 0
		}
		return readline.ENTER
	}
	return readline.CONTINUE
}

func cmdBindKey(L Lua) int {
	key, keyErr := L.ToString(-2)
	if keyErr != nil {
		return L.Push(keyErr)
	}
	key = strings.Replace(strings.ToUpper(key), "-", "_", -1)
	switch L.GetType(-1) {
	case lua.LUA_TFUNCTION:
		chank := L.Dump()
		if err := readline.BindKeyFunc(key, &KeyLuaFuncT{chank}); err != nil {
			return L.Push(nil, err)
		} else {
			return L.Push(true)
		}
	default:
		val, valErr := L.ToString(-1)
		if valErr != nil {
			return L.Push(nil, valErr)
		}
		err := readline.BindKeySymbol(key, val)
		if err != nil {
			return L.Push(nil, err)
		} else {
			return L.Push(true)
		}
	}
}

func cmdGetBindKey(L Lua) int {
	key, keyErr := L.ToString(-1)
	if keyErr != nil {
		return L.Push(nil, keyErr)
	}
	fnc := readline.GetBindKey(key)
	if fnc != nil {
		if stringer, ok := fnc.(fmt.Stringer); ok {
			if str := stringer.String(); str != "" {
				L.PushString(str)
				return 1
			}
		}
	}
	L.PushNil()
	return 1
}
