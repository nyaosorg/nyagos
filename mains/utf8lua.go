package mains

import (
	"strings"

	"github.com/yuin/gopher-lua"
)

func utf8codes(L *lua.LState) int {
	lstr, ok := L.Get(-1).(lua.LString)
	if !ok {
		L.Push(lua.LNil)
		L.Push(lua.LString("invalid utf8"))
		return 2
	}

	p := strings.NewReader(string(lstr))
	pos := 1

	f := func(LL *lua.LState) int {
		r, siz, err := p.ReadRune()
		if err != nil {
			return 0
		}
		LL.Push(lua.LNumber(pos))
		LL.Push(lua.LNumber(r))
		pos += siz
		return 2
	}
	L.Push(L.NewFunction(f))
	L.Push(lstr)
	L.Push(lua.LNumber(1))
	return 3
}

func SetupUtf8Table(L *lua.LState) {
	table := L.NewTable()
	L.SetField(table, "codes", L.NewFunction(utf8codes))
	L.SetField(table, "charpattern", lua.LString("[\000-\x7F\xC2-\xF4][\x80-\xBF]*"))
	L.SetGlobal("utf8", table)
}
