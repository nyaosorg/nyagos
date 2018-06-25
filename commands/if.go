package commands

import (
	"context"
	"errors"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/zetamatta/nyagos/shell"
	"github.com/zetamatta/nyagos/texts"
)

var rxElse = regexp.MustCompile(`(?i)^\s*else`)

func cmdIf(ctx context.Context, cmd Param) (int, error) {
	// if "xxx" == "yyy"
	args := cmd.Args()
	rawargs := cmd.RawArgs()
	not := false
	start := 1

	option := map[string]struct{}{}

	for len(args) >= 2 && strings.HasPrefix(args[1], "/") {
		option[strings.ToLower(args[1])] = struct{}{}
		args = args[1:]
		rawargs = rawargs[1:]
		start++
	}

	if len(args) >= 2 && strings.EqualFold(args[1], "not") {
		not = true
		args = args[1:]
		rawargs = rawargs[1:]
		start++
	}
	status := false
	if len(args) >= 4 && args[2] == "==" {
		if _, ok := option["/i"]; ok {
			status = strings.EqualFold(args[1], args[3])
		} else {
			status = (args[1] == args[3])
		}
		args = args[4:]
		rawargs = rawargs[4:]
		start += 3
	} else if len(args) >= 3 && strings.EqualFold(args[1], "exist") {
		_, err := os.Stat(args[2])
		status = (err == nil)
		args = args[3:]
		rawargs = rawargs[3:]
		start += 2
	} else if len(args) >= 3 && strings.EqualFold(args[1], "errorlevel") {
		num, err := strconv.Atoi(args[2])
		if err == nil {
			status = (shell.LastErrorLevel >= num)
		}
		args = args[3:]
		rawargs = rawargs[3:]
		start += 2
	}

	if not {
		status = !status
	}

	thenBuffer := shell.BufStream{}

	if len(args) > 0 {
		if args[0] == "then" {
			// inline and block `then`
			if len(args) > 1 {
				thenBuffer.Add(strings.Join(rawargs[1:], " "))
			}
			// continue
		} else {
			// inline `then`
			if status {
				return cmd.Spawnlp(ctx, cmd.Args()[start:], cmd.RawArgs()[start:])
			}
			return 0, nil
		}
	}

	// block `then` / `else`

	stream, ok := ctx.Value(shell.StreamID).(shell.Stream)
	if !ok {
		return 1, errors.New("not found stream")
	}

	elseBuffer := shell.BufStream{}
	elsePart := false

	savePrompt := os.Getenv("PROMPT")
	os.Setenv("PROMPT", "if>")
	defer os.Setenv("PROMPT", savePrompt)
	nest := 1
	for {
		_, line, err := cmd.ReadCommand(ctx, stream)
		if err != nil {
			if err != io.EOF {
				return -1, err
			}
			break
		}
		args := texts.SplitLikeShellString(line)
		if len(args) <= 0 {
			continue
		}
		name := strings.ToLower(args[0])
		if _, ok := startList[name]; ok {
			nest++
		} else if name == "end" || name == "endif" {
			nest--
			if nest == 0 {
				break
			}
		} else if name == "else" {
			if nest == 1 {
				elsePart = true
				os.Setenv("PROMPT", "else>")
				line = rxElse.ReplaceAllString(line, "")
			}
		}
		if elsePart {
			elseBuffer.Add(line)
		} else {
			thenBuffer.Add(line)
		}
	}

	if status {
		cmd.Loop(ctx, &thenBuffer)
	} else {
		cmd.Loop(ctx, &elseBuffer)
	}
	return 0, nil
}
