package functions

import (
	"fmt"
	"github.com/zetamatta/go-readline-ny"
)

// CmdGetBindKey is the getter for nyagos.key table.
func CmdGetBindKey(args []any_t) []any_t {
	if len(args) < 1 {
		return []any_t{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[len(args)-1])
	fnc := readline.GlobalKeyMap.GetBindKey(key)
	if fnc != nil {
		return []any_t{fmt.Sprint(fnc)}
	}
	return []any_t{nil}
}
