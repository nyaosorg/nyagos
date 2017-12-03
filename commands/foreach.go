package commands

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/zetamatta/nyagos/shell"
)

var start_list = map[string]bool{
	"foreach": true,
	"if":      true,
}

func cmd_foreach(ctx context.Context, cmd *shell.Cmd) (int, error) {
	stream, ok := ctx.Value("stream").(shell.Stream)

	if !ok {
		return 1, errors.New("Not found stream")
	}

	bufstream := shell.BufStream{}
	save_prompt := os.Getenv("PROMPT")
	os.Setenv("PROMPT", "foreach>")
	nest := 1
	for {
		_, line, err := cmd.ReadCommand(ctx, stream)
		if err != nil {
			break
		}
		args := shell.SplitQ(line)

		name := strings.ToLower(args[0])

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
	if len(cmd.Args) < 2 {
		return 0, nil
	}
	os.Setenv("PROMPT", save_prompt)

	name := cmd.Args[1]
	save := os.Getenv(name)
	for _, value := range cmd.Args[2:] {
		os.Setenv(name, value)
		cmd.Loop(&bufstream)
		bufstream.SetPos(0)
	}
	os.Setenv(name, save)
	return 0, nil
}
