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
	DisableHistory(value bool) bool
}

type _NulStream struct{}

func (stream *_NulStream) ReadLine(ctx context.Context) (context.Context, string, error) {
	return ctx, "", io.EOF
}

func (stream *_NulStream) DisableHistory(value bool) bool {
	return false
}

var NulStream = &_NulStream{}

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
func (sh *Shell) ReadCommand(ctx context.Context) (context.Context, string, error) {
	stream := sh.Stream
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

// Loop executes commands from `stream` until any errors are found.
func (sh *Shell) Loop(ctx0 context.Context, stream Stream) (int, error) {
	backup := sh.Stream
	sh.Stream = stream
	defer func() {
		sh.Stream = backup
	}()

	for {
		ctx, cancel := context.WithCancel(ctx0)

		ctx, line, err := sh.ReadCommand(ctx)
		if err != nil {
			cancel()
			if err == io.EOF {
				return 0, err
			}
			return 1, err
		}

		sigint := make(chan os.Signal)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

		go func(sig chan os.Signal, canc func()) {
			<-sig // wait for receiving signal or channel closed
			canc()
		}(sigint, cancel)

		rc, err := sh.Interpret(ctx, line)
		signal.Stop(sigint)
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
		if rc > 0 {
			fmt.Fprintf(os.Stderr, "exit status %d\n", rc)
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
