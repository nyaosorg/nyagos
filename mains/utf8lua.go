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

func runeCountLua(b []byte) (int, int) {
	if len(b) <= 0 {
		return 0, 0
	}
	if 0x80 <= b[0] && b[0] <= 0xBF && 0xC2 <= b[0] && b[0] <= 0xDF {
		if len(b) < 2 || b[1] < 0x80 || b[1] > 0xBF {
			return -1, 1
		}
	} else if 0xE0 <= b[0] && b[0] <= 0xEF {
		if len(b) < 2 || b[1] < 0x80 || b[1] > 0xBF {
			return -1, 1
		} else if len(b) < 3 || b[2] < 0x80 || b[2] > 0xBF {
			return -1, 2
		}
	} else if 0xF0 <= b[0] && b[0] <= 0xF7 {
		if len(b) < 2 || b[1] < 0x80 || b[1] > 0xBF {
			return -1, 1
		} else if len(b) < 3 || b[2] < 0x80 || b[2] > 0xBF {
			return -1, 2
		} else if len(b) < 4 || b[3] < 0x80 || b[3] > 0xBF {
			return -1, 3
		}
	} else if b[0] >= 0x80 {
		return -1, 0
	}
	return utf8.RuneCount(b), 0
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
	if j <= 0 {
		j += _len
	}
	if i <= 0 {
		i += _len
	}

	pos, errpos := runeCountLua([]byte(s[(i - 1):(j - 1)]))
	if pos < 0 {
		L.Push(lua.LNil)
		L.Push(lua.LNumber(errpos + i))
		return 2
	}
	L.Push(lua.LNumber(pos))
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
