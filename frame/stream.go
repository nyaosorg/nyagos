package frame

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mattn/go-colorable"

	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/readline"
	"github.com/zetamatta/nyagos/shell"
)

type CmdStreamConsole struct {
	shell.CmdSeeker
	DoPrompt func() (int, error)
	History  *history.Container
	Editor   *readline.Editor
	HistPath string
}

var console io.Writer
var prevOptionGoColorable bool = false

var isEscapeSequenceAvailableFlag = false

func GetConsole() io.Writer {
	if isEscapeSequenceAvailableFlag {
		enableVirtualTerminalProcessing()
		console = os.Stdout
	} else if console == nil {
		if isEscapeSequenceAvailable() {
			console = os.Stdout
			enableVirtualTerminalProcessing()
			isEscapeSequenceAvailableFlag = true
		} else {
			console = colorable.NewColorableStdout()
		}
	}
	return console
}

func NewCmdStreamConsole(doPrompt func() (int, error)) *CmdStreamConsole {
	history1 := &history.Container{}
	this := &CmdStreamConsole{
		History: history1,
		Editor: &readline.Editor{
			History: history1,
			Prompt:  doPrompt,
			Writer:  GetConsole()},
		HistPath: filepath.Join(AppDataDir(), "nyagos.history"),
		CmdSeeker: shell.CmdSeeker{
			PlainHistory: []string{},
			Pointer:      -1,
		},
	}
	history1.Load(this.HistPath)
	history1.Save(this.HistPath)
	return this
}

func (this *CmdStreamConsole) ReadLine(ctx context.Context) (context.Context, string, error) {
	if this.Pointer >= 0 {
		if this.Pointer < len(this.PlainHistory) {
			this.Pointer++
			return ctx, this.PlainHistory[this.Pointer-1], nil
		}
		this.Pointer = -1
	}
	var line string
	var err error
	for {
		line, err = this.Editor.ReadLine(ctx)
		if err != nil {
			return ctx, line, err
		}
		var isReplaced bool
		line, isReplaced, err = this.History.Replace(line)
		if err != nil {
			return ctx, line, err
		}
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line != "" {
			break
		}
	}
	row := history.NewHistoryLine(line)
	this.History.PushLine(row)
	fd, err := os.OpenFile(this.HistPath, os.O_APPEND|os.O_CREATE, 0600)
	if err == nil {
		fmt.Fprintln(fd, row.String())
		fd.Close()
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	this.PlainHistory = append(this.PlainHistory, line)
	return ctx, line, err
}
