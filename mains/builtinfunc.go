package mains

import (
	"errors"
	"fmt"

	"github.com/zetamatta/go-box"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/readline"
)

func cstr(value interface{}) (string, bool) {
	if s, ok := value.(string); ok {
		return s, true
	}
	if s, ok := value.(fmt.Stringer); ok {
		return s.String(), true
	}
	if s, ok := value.(int); ok {
		return fmt.Sprintf("%d", s), true
	}
	return "", false
}

func cmdElevated([]interface{}) []interface{} {
	flag, _ := dos.IsElevated()
	return []interface{}{flag}
}

func cmdChdir(args []interface{}) []interface{} {
	if len(args) >= 1 {
		path, ok := cstr(args[0])
		if ok {
			dos.Chdir(path)
			return []interface{}{true}
		}
	}
	return []interface{}{nil, errors.New("directory is required")}
}

func cmdBox(args []interface{}) []interface{} {
	t, ok := args[0].(map[interface{}]interface{})
	if !ok {
		return []interface{}{nil, "Not a table"}
	}
	if len(t) == 0 {
		return []interface{}{}
	}
	sources := make([]string, 0, len(t))
	for _, v := range t {
		if str, ok := cstr(v); ok {
			sources = append(sources, str)
		}
	}
	return []interface{}{box.Choice(sources, readline.Console)}
}
