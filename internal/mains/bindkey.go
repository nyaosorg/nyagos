//go:build !vanilla
// +build !vanilla

package mains

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/gopher-lua"

	"github.com/nyaosorg/go-box/v3"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
	"github.com/nyaosorg/go-readline-ny/nameutils"

	"github.com/nyaosorg/nyagos/internal/texts"
)

type _ReadLineCallBack struct {
	buffer *readline.Buffer
	update func()
}

func (rl *_ReadLineCallBack) Replace(L Lua) int {
	pos, ok := L.Get(-2).(lua.LNumber)
	if !ok {
		return lerror(L, "not a number")
	}
	if pos <= 0 {
		return lerror(L, fmt.Sprintf(":replace: pos=%d: Too small.", pos))
	}
	str := L.ToString(-1)
	posZeroBase := int(pos) - 1
	if posZeroBase > len(rl.buffer.Buffer) {
		return lerror(L, fmt.Sprintf(":replace: pos=%d: Too big.", pos))
	}
	rl.buffer.ReplaceAndRepaint(posZeroBase, string(str))
	L.Push(lua.LTrue)
	L.Push(lua.LNil)
	rl.update()
	return 2
}

func (rl *_ReadLineCallBack) Insert(L Lua) int {
	text := L.ToString(2)
	rl.buffer.InsertAndRepaint(string(text))
	L.Push(lua.LTrue)
	rl.update()
	return 1
}

func (rl *_ReadLineCallBack) evalKey(L Lua) int {
	key, ok := L.Get(-1).(lua.LString)
	if !ok {
		return lerror(L, "eval: expect string as key name or sequence")
	}
	code, ok := keys.NameToCode[keys.NormalizeName(string(key))]
	if !ok {
		code = keys.Code(key)
	}
	function := rl.buffer.LookupCommand(string(code))
	rc := function.Call(L.Context(), rl.buffer)
	rl.buffer.RepaintLastLine()
	switch rc {
	case readline.ENTER:
		L.Push(lua.LTrue)
	case readline.INTR:
		L.Push(lua.LFalse)
	default:
		L.Push(lua.LNil)
	}
	rl.update()
	return 1
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
	case readline.INTR:
		L.Push(lua.LFalse)
	default:
		L.Push(lua.LNil)
	}
	rl.update()
	return 1
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
	box.Println(list, os.Stdout)
	rl.buffer.RepaintAll()
	return 0
}

func (rl *_ReadLineCallBack) Repaint(Lua) int {
	rl.buffer.RepaintLastLine()
	return 0
}

type _KeyLuaFunc struct {
	Chank *lua.LFunction
	L     Lua
}

func (f _KeyLuaFunc) String() string {
	return f.Chank.String()
}

func getPosAndText(b *readline.Buffer) (int, string) {
	pos := -1
	var text strings.Builder
	for i, c := range b.Buffer {
		if i == b.Cursor {
			pos = text.Len() + 1
		}
		c.Moji.WriteTo(&text)
	}
	if pos < 0 {
		pos = text.Len() + 1
	}
	return pos, text.String()
}

func (f *_KeyLuaFunc) Call(ctx context.Context, buffer *readline.Buffer) readline.Result {
	L := f.L
	L.Push(f.Chank)

	pos, text := getPosAndText(buffer)

	rl := &_ReadLineCallBack{buffer: buffer}

	table := L.NewTable()
	L.SetField(table, "pos", lua.LNumber(pos))
	L.SetField(table, "text", lua.LString(text))
	L.SetField(table, "call", L.NewFunction(rl.KeyFunc))
	L.SetField(table, "eval", L.NewFunction(rl.evalKey))
	L.SetField(table, "insert", L.NewFunction(rl.Insert))
	L.SetField(table, "replacefrom", L.NewFunction(rl.Replace))
	L.SetField(table, "lastword", L.NewFunction(rl.LastWord))
	L.SetField(table, "firstword", L.NewFunction(rl.FirstWord))
	L.SetField(table, "boxprint", L.NewFunction(rl.BoxListing))
	L.SetField(table, "repaint", L.NewFunction(rl.Repaint))

	rl.update = func() {
		_pos, _text := getPosAndText(buffer)
		L.SetField(table, "pos", lua.LNumber(_pos))
		L.SetField(table, "text", lua.LString(_text))
	}

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
			if value {
				return readline.ENTER
			} else {
				return readline.INTR
			}
		}
	}
	return readline.CONTINUE
}

func cmdBindKey(L Lua) int {
	key, ok := L.Get(-2).(lua.LString)
	if !ok {
		return lerror(L, "bindkey: key error: "+string(key))
	}
	code, ok := keys.NameToCode[keys.NormalizeName(string(key))]
	if !ok {
		code = keys.Code(key)
	}
	if f, ok := L.Get(-1).(*lua.LFunction); ok {
		readline.GlobalKeyMap.BindKey(code, &_KeyLuaFunc{Chank: f, L: L})
	} else {
		funcname := L.ToString(-1)
		f, ok := readline.NameToFunc[funcname]
		if !ok {
			return lerror(L, "bindkey: func error: "+funcname)
		}
		readline.GlobalKeyMap.BindKey(code, f)
	}
	L.Push(lua.LTrue)
	return 1
}
