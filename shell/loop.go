package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

// Stream is the inteface which can read command-line
type Stream interface {
	ReadLine(context.Context) (context.Context, string, error)
}

func (ses *session) push(lines []string) {
	if lines != nil && len(lines) >= 1 {
		ses.unreadline = append(ses.unreadline, lines...)
	}
}

func (ses *session) pop() (string, bool) {
	if ses.unreadline == nil || len(ses.unreadline) <= 0 {
		return "", false
	}
	line := ses.unreadline[0]
	if len(ses.unreadline) >= 2 {
		ses.unreadline = ses.unreadline[1:]
	} else {
		ses.unreadline = nil
	}
	return line, true
}

// ReadCommand reads completed one command from `stream`.
func (sh *Shell) ReadCommand(ctx context.Context, stream Stream) (context.Context, string, error) {
	var line string
	var err error

	line, ok := sh.pop()
	if !ok {
		ctx, line, err = stream.ReadLine(ctx)
		if err != nil {
			return ctx, line, err
		}

		texts := splitToStatement(line)
		line = texts[0]
		sh.push(texts[1:])
	}
	return ctx, line, nil
}

type streamIDT struct{}

// StreamID is the key-object to find the last stream in the context object.
var StreamID streamIDT

// Loop executes commands from `stream` until any errors are found.
func (sh *Shell) Loop(ctx0 context.Context, stream Stream) (int, error) {
	for {
		ctx, cancel := context.WithCancel(ctx0)
		ctx = context.WithValue(ctx, StreamID, stream)

		ctx, line, err := sh.ReadCommand(ctx, stream)
		if err != nil {
			cancel()
			if err == io.EOF {
				return 0, err
			}
			return 1, err
		}

		sigint := make(chan os.Signal)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		quit := make(chan struct{})

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
		rc, err := sh.Interpret(ctx, line)
		signal.Stop(sigint)
		close(quit)
		close(sigint)

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

// ForEver executes commands from `stream` until EOF are found.
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
