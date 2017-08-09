package mains

import (
	"errors"
	"fmt"

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

func cmdResetCharWidth(args []interface{}) []interface{} {
	readline.ResetCharWidth()
	return []interface{}{}
}

func cmdNetDriveToUNC(args []interface{}) []interface{} {
	if len(args) < 1 {
		return []interface{}{}
	}
	path, ok := args[0].(string)
	if !ok {
		return []interface{}{path}
	}
	unc := dos.NetDriveToUNC(path)
	return []interface{}{unc}
}
