package commands

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/zetamatta/nyagos/shell"
	"github.com/zetamatta/nyagos/texts"
)

var start_list = map[string]bool{
	"foreach": true,
	"if":      true,
}

func cmdForeach(ctx context.Context, cmd Param) (int, error) {
	stream, ok := ctx.Value(shell.StreamId).(shell.Stream)

	if !ok {
		return 1, errors.New("Not found stream")
	}

	bufstream := shell.BufStream{}
	save_prompt := os.Getenv("PROMPT")
	os.Setenv("PROMPT", "foreach>")
	defer os.Setenv("PROMPT", save_prompt)
	nest := 1
	for {
		_, line, err := cmd.ReadCommand(ctx, stream)
		if err != nil {
			if err != io.EOF {
				return -1, err
			}
			break
		}
		name := strings.ToLower(texts.FirstWord(line))
		if _, ok := start_list[name]; ok {
			nest++
		} else if name == "end" {
			nest--
			if nest == 0 {
				break
			}
		}
		bufstream.Add(line)
	}
	if len(cmd.Args()) < 2 {
		return 0, nil
	}

	name := cmd.Arg(1)
	save := os.Getenv(name)
	for _, value := range cmd.Args()[2:] {
		os.Setenv(name, value)
		cmd.Loop(&bufstream)
		bufstream.SetPos(0)
	}
	os.Setenv(name, save)
	return 0, nil
}
