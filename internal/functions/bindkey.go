package functions

import (
	"fmt"
	"github.com/nyaosorg/go-readline-ny/keys"
)

// CmdGetBindKey is the getter for nyagos.key table.
func CmdGetBindKey(param *Param) []any {
	args := param.Args
	if len(args) < 1 {
		return []any{nil, "too few arguments"}
	}
	editor := param.Editor
	if editor == nil {
		return []any{nil, "readline.Editor not found"}
	}
	name := keys.NormalizeName(fmt.Sprint(args[len(args)-1]))
	code, ok := keys.NameToCode[name]
	if !ok {
		code = keys.Code(name)
	}
	command := editor.LookupCommand(string(code))
	if command != nil {
		return []any{command.String()}
	}
	return []any{nil}
}
