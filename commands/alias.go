package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/nyaosorg/nyagos/alias"
)

func cmdAlias(ctx context.Context, cmd Param) (int, error) {
	if len(cmd.Args()) <= 1 {
		for p := alias.Table.Front(); p != nil; p = p.Next() {
			fmt.Fprintf(cmd.Out(), "%s=%s\n", p.Key, p.Value.String())
		}
		return 0, nil
	}
	for _, args := range cmd.Args()[1:] {
		if eqlPos := strings.IndexRune(args, '='); eqlPos >= 0 {
			key := args[0:eqlPos]
			val := args[eqlPos+1:]
			if len(val) > 0 {
				alias.Table.Store(key, alias.New(val))
			} else {
				alias.Table.Delete(key)
			}
		} else {
			val, ok := alias.Table.Load(args)
			if ok {
				fmt.Fprintf(cmd.Out(), "%s=%s\n", args, val.String())
			}
		}
	}
	return 0, nil
}
