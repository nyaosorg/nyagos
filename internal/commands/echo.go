package commands

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
)

func cmdEcho(ctx context.Context, cmd Param) (int, error) {
	io.WriteString(cmd.Out(), strings.Join(cmd.RawArgs()[1:], " "))
	if f, ok := cmd.Out().(*os.File); ok && isatty.IsTerminal(f.Fd()) {
		io.WriteString(f, "\n")
	} else {
		io.WriteString(cmd.Out(), "\r\n")
	}
	return 0, nil
}
