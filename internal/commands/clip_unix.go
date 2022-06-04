//go:build !windows
// +build !windows

package commands

import (
	"context"
	"io"

	"github.com/atotto/clipboard"
)

func cmdClip(ctx context.Context, cmd Param) (int, error) {
	data, err := io.ReadAll(cmd.In())
	if err != nil {
		return 1, err
	}
	clipboard.WriteAll(string(data))
	return 0, nil
}
