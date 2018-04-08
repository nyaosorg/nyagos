package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
)

type Stream interface {
	ReadLine(context.Context) (context.Context, string, error)
}

func (this *session) push(lines []string) {
	if lines != nil && len(lines) >= 1 {
		this.unreadline = append(this.unreadline, lines...)
	}
}

func (this *session) pop() (string, bool) {
	if this.unreadline == nil || len(this.unreadline) <= 0 {
		return "", false
	}
	line := this.unreadline[0]
	if len(this.unreadline) >= 2 {
		this.unreadline = this.unreadline[1:]
	} else {
		this.unreadline = nil
	}
	return line, true
}

func (sh *Shell) ReadCommand(ctx context.Context, stream Stream) (context.Context, string, error) {
	var line string
	var err error

	line, ok := sh.pop()
	if !ok {
		ctx, line, err = stream.ReadLine(ctx)
		if err != nil {
			return ctx, line, err
		}

		texts := SplitToStatement(line)
		line = texts[0]
		sh.push(texts[1:])
	}
	return ctx, line, nil
}

type streamIdT struct{}

var StreamId streamIdT

func (sh *Shell) Loop(ctx0 context.Context, stream Stream) (int, error) {
	sigint := make(chan os.Signal, 1)
	defer close(sigint)
	quit := make(chan struct{}, 1)
	defer close(quit)

	for {
		ctx, cancel := context.WithCancel(ctx0)
		ctx = context.WithValue(ctx, StreamId, stream)

		ctx, line, err := sh.ReadCommand(ctx, stream)
		if err != nil {
			cancel()
			if err == io.EOF {
				return 0, err
			} else {
				return 1, err
			}
		}
		signal.Notify(sigint, os.Interrupt)

		go func(sigint_ chan os.Signal, quit_ chan struct{}, cancel_ func()) {
			for {
				select {
				case <-sigint_:
					cancel_()
					<-quit
					return
				case <-quit:
					cancel_()
					return
				}
			}
		}(sigint, quit, cancel)
		rc, err := sh.InterpretContext(ctx, line)
		signal.Stop(sigint)
		quit <- struct{}{}

		if err != nil {
			if err == io.EOF {
				return rc, err
			}
			if err1, ok := err.(AlreadyReportedError); ok {
				if err1.Err == io.EOF {
					return rc, err
				}
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

func (sh *Shell) ForEver(ctx context.Context, stream Stream) {
	for {
		_, err := sh.Loop(ctx, stream)
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}
