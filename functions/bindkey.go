package functions

import (
	"fmt"
	"github.com/zetamatta/go-readline-ny"
)

// CmdGetBindKey is the getter for nyagos.key table.
func CmdGetBindKey(args []anyT) []anyT {
	if len(args) < 1 {
		return []anyT{nil, "too few arguments"}
	}
	key := fmt.Sprint(args[len(args)-1])
	fnc := readline.GlobalKeyMap.GetBindKey(key)
	if fnc != nil {
		return []anyT{fmt.Sprint(fnc)}
	}
	return []anyT{nil}
}
