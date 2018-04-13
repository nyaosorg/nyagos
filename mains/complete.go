package mains

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/zetamatta/nyagos/completion"
	"github.com/zetamatta/nyagos/lua"
	"github.com/zetamatta/nyagos/readline"
)

var completionHook lua.Object = lua.TNil{}

func luaHookForComplete(ctx context.Context, this *readline.Buffer, rv *completion.List) (*completion.List, error) {
	L, L_ok := ctx.Value(lua.PackageId).(lua.Lua)
	if !L_ok {
		return rv, errors.New("listUpComplete: could not get lua instance")
	}

	L.Push(completionHook)
	if !L.IsFunction(-1) {
		L.Pop(1)
		return rv, nil
	}

	list := make([]string, len(rv.List))
	shownlist := make([]string, len(rv.List))
	for i, v := range rv.List {
		list[i] = v.String()
		shownlist[i] = v.Display()
	}
	L.Push(map[string]interface{}{
		"rawword":   rv.RawWord,
		"pos":       rv.Pos + 1,
		"text":      rv.AllLine,
		"word":      rv.Word,
		"list":      list,
		"shownlist": shownlist,
		"field":     rv.Field,
		"left":      rv.Left,
	})
	if err := L.CallWithContext(ctx, 1, 2); err != nil {
		fmt.Println(err)
		return rv, nil
	}
	if insertStrList, err := L.ToInterface(-2); err == nil {
		if t, ok := insertStrList.(map[interface{}]interface{}); ok {
			listupStrT := t
			if listupStrList, err := L.ToInterface(-1); err == nil {
				if t, ok := listupStrList.(map[interface{}]interface{}); ok {
					listupStrT = t
				}
			}
			list := make([]completion.Element, 0, len(rv.List)+32)
			wordUpr := strings.ToUpper(rv.Word)
			for i := 0; i < len(t); i++ {
				if str, ok := t[i+1].(string); ok {
					strUpr := strings.ToUpper(str)
					if strings.HasPrefix(strUpr, wordUpr) {
						listupStr, ok := listupStrT[i+1].(string)
						if !ok {
							listupStr = str
						}
						list = append(list, completion.Element2{str, listupStr})
					}
				}
			}
			if len(list) > 0 {
				rv.List = list
			}
		}
	}
	L.Pop(2) // remove 2 results.
	return rv, nil
}
