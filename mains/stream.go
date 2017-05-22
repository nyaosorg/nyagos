package mains

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/zetamatta/nyagos/history"
	"github.com/zetamatta/nyagos/readline"
)

type CmdSeeker struct {
	PlainHistory []string
	Pointer      int
}

func (this *CmdSeeker) GetPos() int {
	if this.Pointer >= 0 {
		return this.Pointer
	} else {
		return len(this.PlainHistory)
	}
}

func (this *CmdSeeker) SetPos(pos int) error {
	if pos < len(this.PlainHistory) {
		this.Pointer = pos
		return nil
	} else {
		return errors.New("ICmdStream.SetPos(): Position Overflow")
	}
}

type CmdStreamConsole struct {
	CmdSeeker
	DoPrompt func() (int, error)
	History  *history.Container
	Editor   *readline.Editor
	HistPath string
}

func NewCmdStreamConsole(doPrompt func() (int, error)) *CmdStreamConsole {
	history1 := &history.Container{}
	this := &CmdStreamConsole{
		History:  history1,
		Editor:   &readline.Editor{History: history1, Prompt: doPrompt},
		HistPath: filepath.Join(AppDataDir(), "nyagos.history"),
		CmdSeeker: CmdSeeker{
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
		line, isReplaced = this.History.Replace(line)
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line != "" {
			break
		}
	}
	row := history.NewHistoryLine(line)
	this.History.PushLine(row)
	fd, err := os.OpenFile(this.HistPath, os.O_APPEND, 0600)
	if err != nil && os.IsNotExist(err) {
		fd, err = os.Create(this.HistPath)
	}
	if err == nil {
		fmt.Fprintln(fd, row.String())
		fd.Close()
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	this.PlainHistory = append(this.PlainHistory, line)
	return ctx, line, err
}

type CmdStreamFile struct {
	CmdSeeker
	Scanner *bufio.Scanner
}

func NewCmdStreamFile(r io.Reader) *CmdStreamFile {
	return &CmdStreamFile{
		Scanner: bufio.NewScanner(r),
		CmdSeeker: CmdSeeker{
			PlainHistory: []string{},
			Pointer:      -1,
		},
	}
}

func (this *CmdStreamFile) ReadLine(ctx context.Context) (context.Context, string, error) {
	if this.Pointer >= 0 {
		if this.Pointer < len(this.PlainHistory) {
			this.Pointer++
			return ctx, this.PlainHistory[this.Pointer-1], nil
		}
		this.Pointer = -1
	}
	if !this.Scanner.Scan() {
		if err := this.Scanner.Err(); err != nil {
			return ctx, "", err
		} else {
			return ctx, "", io.EOF
		}
	}
	text := strings.TrimRight(this.Scanner.Text(), "\r\n")
	this.PlainHistory = append(this.PlainHistory, text)
	return ctx, text, nil
}
