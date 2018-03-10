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
	GetPos() int
	SetPos(int) error
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

func (it *Cmd) ReadCommand(ctx context.Context, stream Stream) (context.Context, string, error) {
	var line string
	var err error

	line, ok := it.pop()
	if !ok {
		ctx, line, err = stream.ReadLine(ctx)
		if err != nil {
			return ctx, line, err
		}

		texts := SplitToStatement(line)
		line = texts[0]
		it.push(texts[1:])
	}
	return ctx, line, nil
}

type streamIdT struct{}

var StreamId streamIdT

func (it *Cmd) Loop(stream Stream) (int, error) {
	sigint := make(chan os.Signal, 1)
	defer close(sigint)
	quit := make(chan struct{}, 1)
	defer close(quit)

	var rc int

	for {
		ctx, cancel := context.WithCancel(context.Background())
		ctx = context.WithValue(ctx, StreamId, stream)

		ctx, line, err := it.ReadCommand(ctx, stream)
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
		rc, err = it.InterpretContext(ctx, line)
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
