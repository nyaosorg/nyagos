package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"unsafe"

	"../conio"
	"../lua"
)

type KeyLuaFuncT struct {
	L            lua.Lua
	registoryKey string
}

func getBufferForCallBack(L lua.Lua) (*conio.Buffer, int) {
	if L.GetType(1) != lua.LUA_TTABLE {
		return nil, L.Push(nil, "bindKeyExec: call with : not .")
	}
	L.GetField(1, "buffer")
	if L.GetType(-1) != lua.LUA_TLIGHTUSERDATA {
		return nil, L.Push(nil, "bindKey.Call: invalid object")
	}
	buffer := (*conio.Buffer)(L.ToUserData(-1))
	if buffer == nil {
		return nil, L.Push(nil, "bindKey.Call: invalid member")
	}
	L.Pop(1)
	return buffer, 0
}

func callInsert(L lua.Lua) int {
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

func callKeyFunc(L lua.Lua) int {
	buffer, stackRc := getBufferForCallBack(L)
	if buffer == nil {
		return stackRc
	}
	key, keyErr := L.ToString(2)
	if keyErr != nil {
		return L.Push(nil, keyErr)
	}
	function, funcErr := conio.GetFunc(key)
	if funcErr != nil {
		return L.Push(nil, funcErr)
	}
	switch function.Call(buffer) {
	case conio.ENTER:
		return L.Push(true, true)
	case conio.ABORT:
		return L.Push(true, false)
	default:
		return L.Push(nil)
	}
}

func callLastWord(L lua.Lua) int {
	this, stack_count := getBufferForCallBack(L)
	if this == nil {
		return stack_count
	}
	word, pos := this.CurrentWord()
	return L.Push(word, pos+1)
}

func callFirstWord(L lua.Lua) int {
	this, stack_count := getBufferForCallBack(L)
	if this == nil {
		return stack_count
	}
	word := conio.QuotedFirstWord(this.String())
	return L.Push(word, 0)
}

func callBoxListing(L lua.Lua) int {
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
	conio.BoxPrint(list, os.Stdout)
	this.RepaintAll()
	return 0
}

func (this *KeyLuaFuncT) Call(buffer *conio.Buffer) conio.Result {
	this.L.GetField(lua.LUA_REGISTRYINDEX, this.registoryKey)
	pos := -1
	var text bytes.Buffer
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
	this.L.Push(map[string]interface{}{
		"pos":       pos,
		"text":      text.String(),
		"buffer":    unsafe.Pointer(buffer),
		"call":      callKeyFunc,
		"insert":    callInsert,
		"lastword":  callLastWord,
		"firstword": callFirstWord,
		"boxprint":  callBoxListing,
	})
	if err := this.L.Call(1, 1); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	switch this.L.GetType(-1) {
	case lua.LUA_TSTRING:
		str, strErr := this.L.ToString(-1)
		if strErr == nil {
			buffer.InsertAndRepaint(str)
		}
	case lua.LUA_TBOOLEAN:
		if !this.L.ToBool(-1) {
			buffer.Buffer = []rune{}
			buffer.Length = 0
		}
		return conio.ENTER
	}
	return conio.CONTINUE
}

func cmdBindKey(L lua.Lua) int {
	key, keyErr := L.ToString(-2)
	if keyErr != nil {
		return L.Push(keyErr)
	}
	key = strings.Replace(strings.ToUpper(key), "-", "_", -1)
	switch L.GetType(-1) {
	case lua.LUA_TFUNCTION:
		regkey := "nyagos.bind." + key
		L.SetField(lua.LUA_REGISTRYINDEX, regkey)
		if err := conio.BindKeyFunc(key, &KeyLuaFuncT{L, regkey}); err != nil {
			return L.Push(nil, err)
		} else {
			return L.Push(true)
		}
	default:
		val, valErr := L.ToString(-1)
		if valErr != nil {
			return L.Push(nil, valErr)
		}
		err := conio.BindKeySymbol(key, val)
		if err != nil {
			return L.Push(nil, err)
		} else {
			return L.Push(true)
		}
	}
}
