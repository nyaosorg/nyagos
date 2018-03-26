package commands

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

func cmdEcho(ctx context.Context, cmd Param) (int, error) {
	fmt.Fprint(cmd.Out(), strings.Join(cmd.Args()[1:], " "))
	if f, ok := cmd.Out().(*os.File); ok && isatty.IsTerminal(f.Fd()) {
		fmt.Fprint(f, "\n")
	} else {
		fmt.Fprint(cmd.Out(), "\r\n")
	}
	return 0, nil
}
