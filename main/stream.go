package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"../history"
	"../interpreter"
	"../readline"
)

type ICmdStream interface {
	ReadLine(*context.Context) (string, error)
}

type CmdStreamConsole struct {
	editor   *readline.LineEditor
	history  *history.Container
	histPath string
}

var default_history *history.Container

func NewCmdStreamConsole(it *interpreter.Interpreter) *CmdStreamConsole {
	history1 := new(history.Container)
	editor := readline.NewLineEditor(history1)
	editor.Prompt = printPrompt

	histPath := filepath.Join(AppDataDir(), "nyagos.history")
	history1.Load(histPath)
	history1.Save(histPath)

	default_history = history1

	return &CmdStreamConsole{
		editor:   editor,
		history:  history1,
		histPath: histPath,
	}
}

func (this *CmdStreamConsole) ReadLine(ctx *context.Context) (string, error) {
	*ctx = context.WithValue(*ctx, "history", this.history)
	var line string
	var err error
	for {
		line, err = this.editor.ReadLine(*ctx)
		if err != nil {
			return line, err
		}
		var isReplaced bool
		line, isReplaced = this.history.Replace(line)
		if isReplaced {
			fmt.Fprintln(os.Stdout, line)
		}
		if line != "" {
			break
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}
	row := history.Line{Text: line, Dir: wd, Stamp: time.Now()}
	this.history.PushLine(row)
	fd, err := os.OpenFile(this.histPath, os.O_APPEND, 0600)
	if err != nil && os.IsNotExist(err) {
		// print("create ", this.histPath, "\n")
		fd, err = os.Create(this.histPath)
	}
	if err == nil {
		fmt.Fprintln(fd, row.String())
		fd.Close()
	} else {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	return line, err
}

type CmdStreamFile struct {
	breader *bufio.Reader
}

func NewCmdStreamFile(r io.Reader) *CmdStreamFile {
	return &CmdStreamFile{breader: bufio.NewReader(r)}
}

func (this *CmdStreamFile) ReadLine(ctx *context.Context) (string, error) {
	line, err := this.breader.ReadString('\n')
	if err != nil {
		return "", err
	}
	line = strings.TrimRight(line, "\r\n")
	return line, nil
}

type UnCmdStream struct {
	body  ICmdStream
	queue []string
}

func NewUnCmdStream(body ICmdStream) *UnCmdStream {
	return &UnCmdStream{body: body, queue: nil}
}

func (this *UnCmdStream) ReadLine(ctx *context.Context) (string, error) {
	if this.queue == nil || len(this.queue) <= 0 {
		return this.body.ReadLine(ctx)
	} else {
		line := this.queue[0]
		this.queue = this.queue[1:]
		return line, nil
	}
}

func (this *UnCmdStream) UnreadLine(line string) {
	if this.queue == nil {
		this.queue = []string{line}
	} else {
		this.queue = append(this.queue, line)
	}
}
