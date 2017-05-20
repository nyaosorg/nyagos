package mains

import (
	"errors"
	"fmt"

	"github.com/zetamatta/go-box"

	"../dos"
	"../readline"
)

type Any interface{}

func cstr(value Any) (string, bool) {
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

func cmdElevated([]Any) []Any {
	flag, _ := dos.IsElevated()
	return []Any{flag}
}

func cmdChdir(args []Any) []Any {
	if len(args) >= 1 {
		path, ok := cstr(args[0])
		if ok {
			dos.Chdir(path)
			return []Any{true}
		}
	}
	return []Any{nil, errors.New("directory is required")}
}

func cmdBox(args []Any) []Any {
	t, ok := args[0].(map[Any]Any)
	if !ok {
		return []Any{nil, "Not a table"}
	}
	sources := make([]string, 0, len(t))
	for _, v := range t {
		if str, ok := cstr(v); ok {
			sources = append(sources, str)
		}
	}
	return []Any{box.Choice(sources, readline.Console)}
}
