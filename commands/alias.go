package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/zetamatta/nyagos/alias"
)

func cmdAlias(ctx context.Context, cmd Param) (int, error) {
	if len(cmd.Args()) <= 1 {
		for key, val := range alias.Table {
			fmt.Fprintf(cmd.Out(), "%s=%s\n", key, val.String())
		}
		return 0, nil
	}
	for _, args := range cmd.Args()[1:] {
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
				fmt.Fprintf(cmd.Out(), "%s=%s\n", key, val.String())
			}
		}
	}
	return 0, nil
}
