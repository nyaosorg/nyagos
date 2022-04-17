package shell

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Stream is the inteface which can read command-line
type Stream interface {
	ReadLine(context.Context) (context.Context, string, error)
	DisableHistory(value bool) bool
}

// NulStream is the null implementation for the interface Stream.
type NulStream struct{}

// ReadLine always returns "" and io.EOF.
func (stream *NulStream) ReadLine(ctx context.Context) (context.Context, string, error) {
	return ctx, "", io.EOF
}

// DisableHistory do nothing.
func (stream *NulStream) DisableHistory(value bool) bool {
	return false
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

// endsWithSep returns
//     false when line does not end with `^`
//     true when line ends with `^`
//     false when line ends with `^^`
//     true when line ends with `^^^`
func endsWithSep(line string, contMark byte) bool {
	markCount := 0
	for len(line) > 0 && line[len(line)-1] == contMark {
		markCount++
		line = line[:len(line)-1]
	}
	return markCount%2 != 0
}

// ReadCommand reads completed one command from `stream`.
func (sh *Shell) ReadCommand(ctx context.Context) (context.Context, string, error) {
	stream := sh.Stream
	var line string
	var err error

	line, ok := sh.pop()
	if !ok {
		for {
			outputMutex.Lock()
			os.Stderr.Sync()
			os.Stdout.Sync()
			var line1 string
			ctx, line1, err = stream.ReadLine(ctx)
			outputMutex.Unlock()
			line += line1
			if err != nil {
				return ctx, line, err
			}
			if endsWithSep(line, '^') {
				line = line[:len(line)-1] + "\r\n"
				continue
			}
			if strings.Count(line, "\"")%2 != 0 {
				line = line + "\r\n"
				continue
			}
			break
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

	sigint := make(chan os.Signal)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT)

	defer func() {
		signal.Stop(sigint)
		close(sigint)
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

		go func() {
			select {
			case <-sigint:
				cancel()
			case <-ctx.Done():
			}
		}()

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
