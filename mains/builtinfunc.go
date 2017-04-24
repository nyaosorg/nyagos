package mains

import (
	"errors"
	"fmt"

	"../dos"
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
