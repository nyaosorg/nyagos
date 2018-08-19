// +build !vanilla

package mains

import (
	"strings"
	"unicode/utf8"

	"github.com/yuin/gopher-lua"
)

func utf8char(L *lua.LState) int {
	var buffer strings.Builder
	for i, n := 1, L.GetTop(); i <= n; i++ {
		number, ok := L.Get(i).(lua.LNumber)
		if !ok {
			return lerror(L, "NaN")
		}
		buffer.WriteRune(rune(number))
	}
	L.Push(lua.LString(buffer.String()))
	return 1
}

func utf8codes(L *lua.LState) int {
	lstr, ok := L.Get(-1).(lua.LString)
	if !ok {
		return lerror(L, "invalid utf8")
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

func utf8len(L *lua.LState) int {
	_s, ok := L.Get(1).(lua.LString)
	if !ok {
		return lerror(L, "utf8.len(s,i,j): s is not string")
	}
	s := string(_s)
	_len := len(s)

	j := _len + 1
	i := 1
	top := L.GetTop()
	if top >= 2 {
		if _i, ok := L.Get(2).(lua.LNumber); ok {
			i = int(_i)
		} else {
			return lerror(L, "utf8.len(s,i,j): i is not a number")
		}
		if top >= 3 {
			if _j, ok := L.Get(3).(lua.LNumber); ok {
				j = int(_j)
			} else {
				return lerror(L, "utf8.len(s,i,j): j is not a number)")
			}
		}
	}
	if j < 0 {
		j += _len + 1
	} else if j > 0 {
		j--
	}
	if i < 0 {
		i += _len + 1
	} else if i > 0 {
		i--
	}
	if !utf8.RuneStart(s[i]) {
		return lerror(L, "utf8.len: not start byte")
	}
	s = s[i:]
	j -= i
	length := 0
	for pos := range s {
		if pos > j {
			break
		}
		length++
	}
	L.Push(lua.LNumber(length))
	return 1
}

func utf8offset(L *lua.LState) int {
	_s, ok := L.Get(1).(lua.LString)
	if !ok {
		return lerror(L, "utf8.offset: not string")
	}
	s := string(_s)
	_n, ok := L.Get(2).(lua.LNumber)
	if !ok {
		return lerror(L, "utf8.offset: not a number")
	}
	n := int(_n)

	_i, ok := L.Get(3).(lua.LNumber)
	i := 1
	if ok {
		i = int(_i)
	}
	if i == 0 {
		i = 1
	} else if i < 0 {
		i += len(s) + 1
	}
	for i > 0 && i < len(s) && !utf8.RuneStart(s[i-1]) {
		i++
	}
	s = s[i-1:]

	if n < 0 {
		n = utf8.RuneCountInString(s) + n
	} else if n > 0 {
		n--
	}
	for pos := range s {
		if n <= 0 {
			L.Push(lua.LNumber(pos + 1 + i - 1))
			return 1
		}
		n--
	}
	L.Push(lua.LNumber(len(s)))
	return 1
}

func setupUtf8Table(L *lua.LState) {
	table := L.NewTable()
	L.SetField(table, "codes", L.NewFunction(utf8codes))
	L.SetField(table, "charpattern", lua.LString("[\000-\x7F\xC2-\xF4][\x80-\xBF]*"))
	L.SetField(table, "char", L.NewFunction(utf8char))
	L.SetField(table, "len", L.NewFunction(utf8len))
	L.SetField(table, "offset", L.NewFunction(utf8offset))
	L.SetGlobal("utf8", table)
}
