package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/zetamatta/nyagos/shell"
)

func cmd_chmod_(args []string) error {
	if len(args) < 2 {
		return errors.New("Usage: chmod ooo (files...)")
	}
	permission := 0
	if len(args[0]) != 3 {
		return fmt.Errorf("%s: invalid permission str", args[0])
	}
	for _, r := range args[0] {
		n := strings.IndexRune("01234567", r)
		if n < 0 {
			return fmt.Errorf("%s: invalid permission str", args[0])
		}
		permission = permission*8 + n
	}
	for _, fname := range args[1:] {
		if err := os.Chmod(fname, os.FileMode(permission)); err != nil {
			return fmt.Errorf("%s: %s", fname, err.Error())
		}
	}
	return nil
}

func cmd_chmod(_ context.Context, cmd *shell.Cmd) (int, error) {
	if err := cmd_chmod_(cmd.Args[1:]); err != nil {
		return 1, err
	} else {
		return 0, nil
	}

}
