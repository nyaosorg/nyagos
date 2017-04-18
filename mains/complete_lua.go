package mains

import (
	"errors"
	"fmt"
	"strings"

	"../completion"
	"../lua"
	"../readline"
)

var completionHook lua.Pushable = lua.TNil{}

func luaHookForComplete(this *readline.Buffer, rv *completion.List) (*completion.List, error) {
	L, L_ok := this.Context.Value("lua").(lua.Lua)
	if !L_ok {
		return rv, errors.New("listUpComplete: could not get lua instance")
	}

	L.Push(completionHook)
	if L.IsFunction(-1) {
		L.Push(map[string]interface{}{
			"rawword": rv.RawWord,
			"pos":     rv.Pos + 1,
			"text":    rv.AllLine,
			"word":    rv.Word,
		})
		L.NewTable()
		for key, val := range rv.List {
			L.Push(1 + key)
			L.PushString(val.InsertStr)
			L.SetTable(-3)
		}
		L.SetField(-2, "list")
		if err := L.Call(1, 1); err != nil {
			fmt.Println(err)
		}
		if L.IsTable(-1) {
			list := make([]completion.Element, 0, len(rv.List)+32)
			wordUpr := strings.ToUpper(rv.Word)
			for i := 1; true; i++ {
				L.Push(i)
				L.GetTable(-2)
				str, strErr := L.ToString(-1)
				L.Pop(1)
				if strErr != nil || str == "" {
					break
				}
				strUpr := strings.ToUpper(str)
				if strings.HasPrefix(strUpr, wordUpr) {
					list = append(list, completion.Element{str, str})
				}
			}
			if len(list) > 0 {
				rv.List = list
			}
		}
	}
	L.Pop(1) // remove something not function or result-table
	return rv, nil
}
