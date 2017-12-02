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

func (it *Cmd) ReadCommand(ctx context.Context, stream Stream) (context.Context, string, error) {
	var line string
	var err error

	if it.Unreadline != nil && len(it.Unreadline) > 0 {
		line = it.Unreadline[0]
		if len(it.Unreadline) >= 2 {
			it.Unreadline = it.Unreadline[1:]
		} else {
			it.Unreadline = nil
		}
	} else {
		ctx, line, err = stream.ReadLine(ctx)
		if err != nil {
			return ctx, line, err
		}
		texts := SplitToStatement(line)
		line = texts[0]
		if len(texts) >= 2 {
			it.Unreadline = texts[1:]
		} else {
			it.Unreadline = nil
		}
	}
	return ctx, line, nil
}

func (it *Cmd) Loop(stream Stream) error {
	sigint := make(chan os.Signal, 1)
	defer close(sigint)
	quit := make(chan struct{}, 1)
	defer close(quit)

	for {
		ctx, cancel := context.WithCancel(context.Background())
		ctx = context.WithValue(ctx, "stream", stream)

		ctx, line, err := it.ReadCommand(ctx, stream)
		if err != nil {
			cancel()
			return err
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
		_, err = it.InterpretContext(ctx, line)
		signal.Stop(sigint)
		quit <- struct{}{}

		if err != nil {
			if err == io.EOF {
				break
			}
			if err1, ok := err.(AlreadyReportedError); ok {
				if err1.Err == io.EOF {
					break
				}
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}
	return nil
}
