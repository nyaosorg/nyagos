package commands

import (
	"context"
	"io"
	"unicode/utf8"

	"github.com/atotto/clipboard"
	"github.com/nyaosorg/go-windows-mbcs"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

func cmdClip(ctx context.Context, cmd Param) (int, error) {
	if isTerminalIn(cmd.In()) {
		c, err := nodos.EnableProcessInput()
		if err != nil {
			return -1, err
		}
		defer c()
	}
	data, err := io.ReadAll(cmd.In())
	if err != nil {
		return 1, err
	}
	if utf8.Valid(data) {
		clipboard.WriteAll(string(data))
	} else {
		str, err := mbcs.AnsiToUtf8(data, mbcs.ConsoleCP())
		if err == nil {
			clipboard.WriteAll(str)
		} else {
			return 2, err
		}
	}
	return 0, nil
}
