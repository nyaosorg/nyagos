package mains

import (
	"errors"
	"fmt"

	"github.com/zetamatta/go-box"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/readline"
)

type any_t = interface{}

func cmdElevated([]any_t) []any_t {
	flag, _ := dos.IsElevated()
	return []any_t{flag}
}

func cmdChdir(args []any_t) []any_t {
	if len(args) >= 1 {
		dos.Chdir(fmt.Sprint(args[0]))
		return []any_t{true}
	}
	return []any_t{nil, errors.New("directory is required")}
}

func cmdBox(args []any_t) []any_t {
	t, ok := args[0].(map[any_t]any_t)
	if !ok {
		return []any_t{nil, "Not a table"}
	}
	if len(t) == 0 {
		return []any_t{}
	}
	sources := make([]string, 0, len(t))
	for i, i_ := 1, len(t); i <= i_; i++ {
		if val, ok := t[i]; ok {
			sources = append(sources, fmt.Sprint(val))
		}
	}
	return []any_t{box.Choice(sources, readline.Console)}
}

func cmdResetCharWidth(args []any_t) []any_t {
	readline.ResetCharWidth()
	return []any_t{}
}

func cmdNetDriveToUNC(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{}
	}
	path, ok := args[0].(string)
	if !ok {
		return []any_t{path}
	}
	unc := dos.NetDriveToUNC(path)
	return []any_t{unc}
}
