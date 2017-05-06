package mains

import (
	"errors"
	"fmt"
	"strings"

	"../completion"
	"../lua"
	"../readline"
)

var completionHook lua.Object = lua.TNil{}

func luaHookForComplete(this *readline.Buffer, rv *completion.List) (*completion.List, error) {
	L, L_ok := this.Context.Value("lua").(lua.Lua)
	if !L_ok {
		return rv, errors.New("listUpComplete: could not get lua instance")
	}

	L.Push(completionHook)
	if L.IsFunction(-1) {
		list := make([]string, len(rv.List))
		for i, v := range rv.List {
			list[i] = v.InsertStr
		}
		L.Push(map[string]interface{}{
			"rawword": rv.RawWord,
			"pos":     rv.Pos + 1,
			"text":    rv.AllLine,
			"word":    rv.Word,
			"list":    list,
		})
		if err := L.Call(1, 1); err != nil {
			fmt.Println(err)
		}
		result, err := L.ToInterface(-1)
		if err == nil {
			if t, ok := result.(map[interface{}]interface{}); ok {
				list := make([]completion.Element, 0, len(rv.List)+32)
				wordUpr := strings.ToUpper(rv.Word)
				for _, v := range t {
					str, ok := v.(string)
					if ok {
						strUpr := strings.ToUpper(str)
						if strings.HasPrefix(strUpr, wordUpr) {
							list = append(list, completion.Element{str, str})
						}
					}
				}
				if len(list) > 0 {
					rv.List = list
				}
			}
		}
	}
	L.Pop(1) // remove something not function or result-table
	return rv, nil
}
