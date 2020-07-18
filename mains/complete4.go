// +build !vanilla

package mains

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/gopher-lua"

	"github.com/zetamatta/nyagos/completion"
)

func complete4getter(L Lua) int {
	key, ok := L.Get(-1).(lua.LString)
	if !ok {
		return lerror(L, "nyagos.complete_for[] too few arguments")
	}
	if p, ok := completion.CustomCompletion[string(key)]; ok {
		if c, ok := p.(*customCompleter); ok {
			L.Push(c.Func)
		} else {
			ud := L.NewUserData()
			ud.Value = p
			L.Push(ud)
		}
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

type customCompleter struct {
	Func *lua.LFunction
	Name string
}

func (c *customCompleter) String() string {
	return c.Name
}

func (c *customCompleter) Complete(ctx context.Context, ua completion.UncCompletion, args []string) ([]completion.Element, error) {
	LL, ok := ctx.Value(luaKey).(Lua)
	if !ok {
		return nil, errors.New("completion.CustomCompletion: no lua instance")
	}
	tbl := LL.NewTable()
	for i, arg1 := range args {
		LL.SetTable(tbl, lua.LNumber(i+1), lua.LString(arg1))
	}

	defer setContext(LL, getContext(LL))
	setContext(LL, ctx)

	LL.Push(c.Func)
	LL.Push(tbl)
	if err := LL.PCall(1, 1, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	result := LL.Get(-1)
	if rtbl, ok := result.(*lua.LTable); ok {
		r := make([]completion.Element, 0, rtbl.Len())
		var base string
		if len(args) > 0 {
			base = strings.ToUpper(args[len(args)-1])
		}
		rtbl.ForEach(func(_ lua.LValue, val lua.LValue) {
			if _s, ok := val.(lua.LString); ok {
				s := string(_s)
				if strings.HasPrefix(strings.ToUpper(s), base) {
					r = append(r, completion.Element1(s))
				}
			}
		})
		return r, nil
	} else if s, ok := result.(lua.LString); ok {
		list := strings.Split(string(s), "\n")
		r := make([]completion.Element, 0, len(list))
		for _, r1 := range list {
			r = append(r, completion.Element1(string(r1)))
		}
		return r, nil
	} else {
		return nil, errors.New("not a table or string")
	}
}

func complete4setter(L Lua) int {
	key, ok := L.Get(-2).(lua.LString)
	if !ok {
		return lerror(L, "nyagos.complete_for[] too few arguments")
	}
	val := L.Get(-1)
	if val == lua.LNil {
		delete(completion.CustomCompletion, string(key))
		L.Push(lua.LTrue)
		return 1
	}
	if f, ok := val.(*lua.LFunction); ok {
		completion.CustomCompletion[string(key)] = &customCompleter{
			Func: f,
			Name: string(key),
		}
		L.Push(lua.LTrue)
		return 1
	}
	if ud, ok := val.(*lua.LUserData); ok {
		if c, ok := ud.Value.(completion.CustomCompleter); ok {
			completion.CustomCompletion[string(key)] = c
			L.Push(lua.LTrue)
			return 1
		}
	}
	return lerror(L, "nyagos.complete_for[]= not function")
}
