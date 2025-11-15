package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/nyaosorg/go-readline-ny"
)

type Editor interface {
	StoreDefault(string)
}

// Stream is the inteface which can read command-line
type Stream interface {
	ReadLine(context.Context) (string, error)
	GetHistory() History
	GetEditor() Editor
}

// NulStream is the null implementation for the interface Stream.
type NulStream struct {
	super Stream
}

// ReadLine always returns "" and io.EOF.
func (stream *NulStream) ReadLine(ctx context.Context) (string, error) {
	return "", io.EOF
}

func (stream *NulStream) GetHistory() History {
	if stream.super != nil {
		return stream.super.GetHistory()
	}
	return &_NulHistory{}
}

func (stream *NulStream) GetEditor() Editor {
	if stream.super != nil {
		return stream.super.GetEditor()
	}
	return nil
}

func (ses *session) push(lines []string) {
	if len(lines) >= 1 {
		ses.unreadline = append(ses.unreadline, lines...)
	}
}

func (ses *session) pop() (string, bool) {
	if len(ses.unreadline) <= 0 {
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
func (sh *Shell) ReadCommand(ctx context.Context) (string, error) {
	stream := sh.Stream
	var line string
	var err error

	line, ok := sh.pop()
	if !ok {
		outputMutex.Lock()
		os.Stderr.Sync()
		os.Stdout.Sync()
		line, err = stream.ReadLine(ctx)
		outputMutex.Unlock()
		if err != nil {
			return line, err
		}

		texts := splitToStatement(line)
		line = texts[0]
		sh.push(texts[1:])
	}
	return line, nil
}

func dropSigint(sigint chan os.Signal) {
	for {
		select {
		case <-sigint:
			// println("drop sigint")
			runtime.Gosched()
		default:
			return
		}
	}
}

// Loop executes commands from `stream` until any errors are found.
func (sh *Shell) Loop(ctx0 context.Context, stream Stream) (int, error) {
	backup := sh.Stream
	sh.Stream = stream
	defer func() {
		sh.Stream = backup
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT)

	defer func() {
		signal.Stop(sigint)
		close(sigint)
	}()

	var cancel context.CancelFunc

	go func() {
		for range sigint {
			if cancel != nil {
				cancel()
			}
		}
	}()

	for {
		dropSigint(sigint)

		var ctx context.Context
		ctx, cancel = context.WithCancel(ctx0)

		line, err := sh.ReadCommand(ctx)
		if err != nil {
			cancel()
			cancel = nil
			if err == io.EOF {
				return 0, err
			}
			if err == readline.CtrlC {
				fmt.Fprintln(os.Stderr, err.Error())
				continue
			}
			return 1, err
		}

		rc, err := sh.Interpret(ctx, line)

		if err != nil {
			if err == io.EOF {
				cancel()
				return rc, err
			}
			if err1, ok := err.(AlreadyReportedError); ok {
				if err1.Err == io.EOF {
					cancel()
					return rc, err
				}
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		if rc > 0 {
			fmt.Fprintf(os.Stderr, "exit status %d\n", rc)
		}
		cancel()
		cancel = nil
	}
}

// ForEver executes commands from `stream` until EOF are found.
func (sh *Shell) ForEver(ctx context.Context, stream Stream) error {
	for {
		_, err := sh.Loop(ctx, stream)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}
