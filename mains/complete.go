//go:build !vanilla
// +build !vanilla

package mains

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yuin/gopher-lua"

	"github.com/zetamatta/go-readline-ny"

	"github.com/zetamatta/nyagos/completion"
)

func luaHookForComplete(ctx context.Context, this *readline.Buffer, rv *completion.List) (*completion.List, error) {
	L, ok := ctx.Value(luaKey).(Lua)
	if !ok {
		return rv, errors.New("listUpComplete: could not get lua instance")
	}

	nyagosTbl, ok := L.GetGlobal("nyagos").(*lua.LTable)
	if !ok {
		return rv, nil
	}
	f, ok := L.GetField(nyagosTbl, "completion_hook").(*lua.LFunction)
	if !ok {
		return rv, nil
	}

	list := L.NewTable()
	shownlist := L.NewTable()
	for i, v := range rv.List {
		L.SetTable(list, lua.LNumber(i+1), lua.LString(v.String()))
		L.SetTable(shownlist, lua.LNumber(i+1), lua.LString(v.Display()))
	}
	tbl := L.NewTable()
	L.SetField(tbl, "rawword", lua.LString(rv.RawWord))
	L.SetField(tbl, "pos", lua.LNumber(rv.Pos+1))
	L.SetField(tbl, "text", lua.LString(rv.AllLine))
	L.SetField(tbl, "word", lua.LString(rv.Word))
	L.SetField(tbl, "list", list)
	L.SetField(tbl, "shownlist", shownlist)
	field := L.NewTable()
	for key, val := range rv.Field {
		L.SetTable(field, lua.LNumber(key+1), lua.LString(val))
	}
	L.SetField(tbl, "field", field)
	L.SetField(tbl, "left", lua.LString(rv.Left))

	defer setContext(getContext(L), L)
	setContext(ctx, L)

	L.Push(f)
	L.Push(tbl)

	if err := L.PCall(1, 2, nil); err != nil {
		fmt.Println(err)
		return rv, nil
	}

	defer L.Pop(2) // remove 2 results.

	insertStrs, ok := L.Get(-2).(*lua.LTable)
	if !ok {
		return rv, nil
	}
	listupStrs, ok := L.Get(-1).(*lua.LTable)
	if !ok {
		listupStrs = insertStrs
	}
	newList := make([]completion.Element, 0, len(rv.List)+32)
	wordUpr := strings.ToUpper(rv.Word)
	L.ForEach(insertStrs, func(key, val lua.LValue) {
		str, ok := val.(lua.LString)
		if ok {
			strUpr := strings.ToUpper(string(str))
			if strings.HasPrefix(strUpr, wordUpr) {
				listupStr, ok := L.GetTable(listupStrs, key).(lua.LString)
				if !ok {
					listupStr = str
				}
				newList = append(newList, completion.Element2{
					string(str), string(listupStr)})
			}
		}
	})
	if len(newList) > 0 {
		rv.List = newList
	}
	return rv, nil
}
