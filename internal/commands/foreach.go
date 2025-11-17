package commands

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/nyaosorg/nyagos/internal/shell"
	"github.com/nyaosorg/nyagos/internal/texts"
)

var startList = map[string]bool{
	"foreach": true,
	"if":      true,
}

func cmdForeach(ctx context.Context, cmd Param) (int, error) {
	bufstream := &shell.BufStream{
		History: cmd.GetHistory(),
	}
	savePrompt := os.Getenv("PROMPT")
	os.Setenv("PROMPT", "foreach>")
	defer os.Setenv("PROMPT", savePrompt)
	nest := 1
	for {
		line, err := cmd.ReadCommand(ctx)
		if err != nil {
			if err != io.EOF {
				return -1, err
			}
			break
		}
		name := strings.ToLower(texts.FirstWord(line))
		if _, ok := startList[name]; ok {
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
		cmd.Loop(ctx, bufstream)
		bufstream.SetPos(0)
	}
	os.Setenv(name, save)
	return 0, nil
}
