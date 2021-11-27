package frame

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-readline-ny"
	"github.com/nyaosorg/go-windows-consoleicon"

	"github.com/nyaosorg/nyagos/history"
	"github.com/nyaosorg/nyagos/shell"
)

type CmdStreamConsole struct {
	shell.CmdSeeker
	DoPrompt func() (int, error)
	History  *history.Container
	Editor   *readline.Editor
	HistPath string
}

func NewCmdStreamConsole(doPrompt func() (int, error)) *CmdStreamConsole {
	history1 := &history.Container{}
	stream := &CmdStreamConsole{
		History: history1,
		Editor: &readline.Editor{
			History:  history1,
			Prompt:   doPrompt,
			Writer:   colorable.NewColorableStdout(),
			Coloring: &_Coloring{}},
		HistPath: filepath.Join(appDataDir(), "nyagos.history"),
		CmdSeeker: shell.CmdSeeker{
			PlainHistory: []string{},
			Pointer:      -1,
		},
	}
	history1.Load(stream.HistPath)
	history1.Save(stream.HistPath)
	return stream
}

func (stream *CmdStreamConsole) DisableHistory(value bool) bool {
	return stream.History.IgnorePush(value)
}

func (stream *CmdStreamConsole) ReadLine(ctx context.Context) (context.Context, string, error) {
	if stream.Pointer >= 0 {
		if stream.Pointer < len(stream.PlainHistory) {
			stream.Pointer++
			return ctx, stream.PlainHistory[stream.Pointer-1], nil
		}
		stream.Pointer = -1
	}
	var line string
	var err error
	for {
		disabler := colorable.EnableColorsStdout(nil)
		clean, err2 := consoleicon.SetFromExe()
		for {
			line, err = stream.Editor.ReadLine(ctx)
			if err != readline.CtrlC {
				break
			}
			fmt.Fprintln(os.Stderr, err.Error())
		}
		if err2 == nil {
			clean(false)
		}
		disabler()
		if err != nil {
			return ctx, line, err
		}
		var isReplaced bool
		line, isReplaced, err = stream.History.Replace(line)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line != "" {
			break
		}
	}
	row := history.NewHistoryLine(line)
	stream.History.PushLine(row)
	fd, err := os.OpenFile(stream.HistPath, os.O_APPEND|os.O_CREATE, 0600)
	if err == nil {
		fmt.Fprintln(fd, row.String())
		fd.Close()
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	stream.PlainHistory = append(stream.PlainHistory, line)
	return ctx, line, err
}
