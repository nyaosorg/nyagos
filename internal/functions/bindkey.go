package functions

import (
	"fmt"
	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-readline-ny/keys"
)

// CmdGetBindKey is the getter for nyagos.key table.
func (*Env) CmdGetBindKey(args []any) []any {
	if len(args) < 1 {
		return []any{nil, "too few arguments"}
	}
	name := keys.NormalizeName(fmt.Sprint(args[len(args)-1]))
	code, ok := keys.NameToCode[name]
	if !ok {
		code = keys.Code(name)
	}
	command, ok := readline.GlobalKeyMap.Lookup(code)
	if ok {
		return []any{command.String()}
	}
	return []any{nil}
}
