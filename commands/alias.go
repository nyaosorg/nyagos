package commands

import "fmt"
import "strings"

import "../alias"
import "../interpreter"

func cmd_alias(cmd *interpreter.Interpreter) (interpreter.NextT, error) {
	if len(cmd.Args) <= 1 {
		for key, val := range alias.Table {
			fmt.Fprintf(cmd.Stdout, "%s=%s\n", key, val.String())
		}
		return interpreter.CONTINUE, nil
	}
	for _, args := range cmd.Args[1:] {
		if eqlPos := strings.IndexRune(args, '='); eqlPos >= 0 {
			key := args[0:eqlPos]
			val := args[eqlPos+1:]
			if len(val) > 0 {
				alias.Table[strings.ToLower(key)] = alias.New(val)
			} else {
				delete(alias.Table, strings.ToLower(key))
			}
		} else {
			key := strings.ToLower(args)
			val, ok := alias.Table[key]
			if ok {
				fmt.Fprintf(cmd.Stdout, "%s=%s\n", key, val.String())
			}
		}
	}
	return interpreter.CONTINUE, nil
}
