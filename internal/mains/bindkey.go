//go:build !vanilla
// +build !vanilla

package mains

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/gopher-lua"

	"github.com/nyaosorg/go-box/v2"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/nameutils"

	"github.com/nyaosorg/nyagos/internal/texts"
)

type _KeyLuaFunc struct {
	Chank *lua.LFunction
}

type _ReadLineCallBack struct {
	buffer *readline.Buffer
}

func (rl *_ReadLineCallBack) Replace(L Lua) int {
	pos, ok := L.Get(-2).(lua.LNumber)
	if !ok {
		return lerror(L, "not a number")
	}
	str := L.ToString(-1)
	posZeroBase := int(pos) - 1
	if posZeroBase > len(rl.buffer.Buffer) {
		return lerror(L, fmt.Sprintf(":replace: pos=%d: Too big.", pos))
	}
	rl.buffer.ReplaceAndRepaint(posZeroBase, string(str))
	L.Push(lua.LTrue)
	L.Push(lua.LNil)
	return 2
}

func (rl *_ReadLineCallBack) Insert(L Lua) int {
	text := L.ToString(2)
	rl.buffer.InsertAndRepaint(string(text))
	L.Push(lua.LTrue)
	return 1
}

func (rl *_ReadLineCallBack) evalKey(L Lua) int {
	key := L.ToString(2)
	function := rl.buffer.LookupCommand(key)
	rc := function.Call(L.Context(), rl.buffer)
	rl.buffer.RepaintLastLine()
	switch rc {
	case readline.ENTER:
		L.Push(lua.LTrue)
		L.Push(lua.LTrue)
		return 2
	case readline.ABORT:
		L.Push(lua.LTrue)
		L.Push(lua.LFalse)
		return 2
	default:
		L.Push(lua.LNil)
		return 1
	}
}

func (rl *_ReadLineCallBack) KeyFunc(L Lua) int {
	key := L.ToString(2)
	function, err := nameutils.GetFunc(key)
	if err != nil {
		return lerror(L, err.Error())
	}
	switch function.Call(L.Context(), rl.buffer) {
	case readline.ENTER:
		L.Push(lua.LTrue)
		L.Push(lua.LTrue)
		return 2
	case readline.ABORT:
		L.Push(lua.LTrue)
		L.Push(lua.LFalse)
		return 2
	default:
		L.Push(lua.LNil)
		return 1
	}
}

func (rl *_ReadLineCallBack) LastWord(L Lua) int {
	word, pos := rl.buffer.CurrentWord()
	L.Push(lua.LString(word))
	L.Push(lua.LNumber(pos + 1))
	return 2
}

func (rl *_ReadLineCallBack) FirstWord(L Lua) int {
	word := texts.FirstWord(rl.buffer.String())
	L.Push(lua.LString(word))
	L.Push(lua.LNumber(0))
	return 2
}

func (rl *_ReadLineCallBack) BoxListing(L Lua) int {
	// stack +1: readline.Buffer
	// stack +2: table
	// stack +3: index or value
	fmt.Print("\n")
	table := L.ToTable(2)
	size := table.Len()
	list := make([]string, size)
	for i := 0; i < size; i++ {
		list[i] = L.GetTable(table, lua.LNumber(i+1)).String()
	}
	box.Print(L.Context(), list, os.Stdout)
	rl.buffer.RepaintAll()
	return 0
}

func (f _KeyLuaFunc) String() string {
	return f.Chank.String()
}
func (f *_KeyLuaFunc) Call(ctx context.Context, buffer *readline.Buffer) readline.Result {
	L, ok := ctx.Value(luaKey).(Lua)
	if !ok {
		println("(*mains._KeyLuaFunc)Call: lua instance not found")
		return readline.CONTINUE
	}
	L.Push(f.Chank)
	pos := -1
	var text strings.Builder
	for i, c := range buffer.Buffer {
		if i == buffer.Cursor {
			pos = text.Len() + 1
		}
		c.Moji.WriteTo(&text)
	}
	if pos < 0 {
		pos = text.Len() + 1
	}

	rl := &_ReadLineCallBack{buffer: buffer}

	table := L.NewTable()
	L.SetField(table, "pos", lua.LNumber(pos))
	L.SetField(table, "text", lua.LString(text.String()))
	L.SetField(table, "call", L.NewFunction(rl.KeyFunc))
	L.SetField(table, "eval", L.NewFunction(rl.evalKey))
	L.SetField(table, "insert", L.NewFunction(rl.Insert))
	L.SetField(table, "replacefrom", L.NewFunction(rl.Replace))
	L.SetField(table, "lastword", L.NewFunction(rl.LastWord))
	L.SetField(table, "firstword", L.NewFunction(rl.FirstWord))
	L.SetField(table, "boxprint", L.NewFunction(rl.BoxListing))

	defer setContext(getContext(L), L)
	setContext(ctx, L)

	L.Push(table)
	err := L.PCall(1, 1, nil)
	if err != nil {
		println(err.Error())
	} else {
		switch value := L.Get(-1).(type) {
		case lua.LString:
			buffer.InsertAndRepaint(string(value))
		case lua.LBool:
			if !value {
				buffer.Buffer = buffer.Buffer[:0]
			}
			return readline.ENTER
		}
	}
	return readline.CONTINUE
}

func cmdBindKey(L Lua) int {
	keyTmp, ok := L.Get(-2).(lua.LString)
	if !ok {
		return lerror(L, "bindkey: key error")
	}
	key := strings.Replace(strings.ToUpper(string(keyTmp)), "-", "_", -1)
	switch value := L.Get(-1).(type) {
	case *lua.LFunction:
		if err := nameutils.BindKeyFunc(readline.GlobalKeyMap, key, &_KeyLuaFunc{value}); err != nil {
			return lerror(L, err.Error())
		}
		L.Push(lua.LTrue)
		return 1
	default:
		val := L.ToString(-1)
		err := nameutils.BindKeySymbol(readline.GlobalKeyMap, key, val)
		if err != nil {
			return lerror(L, err.Error())
		}
		L.Push(lua.LTrue)
		return 1
	}
}
