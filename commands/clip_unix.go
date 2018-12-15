// +build !windows

package commands

import (
	"context"
	"io/ioutil"

	"github.com/atotto/clipboard"
)

func cmdClip(ctx context.Context, cmd Param) (int, error) {
	data, err := ioutil.ReadAll(cmd.In())
	if err != nil {
		return 1, err
	}
	clipboard.WriteAll(string(data))
	return 0, nil
}
