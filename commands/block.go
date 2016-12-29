package commands

import (
	"context"
	"errors"
	"os/exec"
	"strings"

	"../text"
)

type UnreadLine interface {
	ReadLine(*context.Context) (string, error)
	UnreadLine(string)
}

type IInterpreter interface {
	InterpretContext(ctx context.Context, line string) (int, error)
	GetRawArgs() []string
}

func read_block(ctx context.Context, run bool, start int) (rc int, err error) {
	readline, readline_ok := ctx.Value("readline").(UnreadLine)
	if !readline_ok {
		return -1, errors.New("context 'readline' not found")
	}
	shell, shell_ok := ctx.Value("interpreter").(IInterpreter)
	if !shell_ok {
		return -1, errors.New("context 'interpreter' not found")
	}
	rawargs := shell.GetRawArgs()
	if start < len(rawargs) {
		readline.UnreadLine(strings.Join(rawargs[start:], " "))
	}
	lines, err := text.ReadBlock(func() (string, error) {
		return readline.ReadLine(&ctx)
	}, func(line string) {
		readline.UnreadLine(line)
	}), nil

	if run {
		for _, line := range lines {
			rc, err = shell.InterpretContext(ctx, line)
		}
	}
	return
}

func cmd_block(ctx context.Context, cmd *exec.Cmd) (rc int, err error) {
	return read_block(ctx, true, 1)
}
