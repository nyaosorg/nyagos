package shell

import (
	"bufio"
	"context"
	"io"
	"os"
	"strings"
)

type CmdSeeker struct {
	PlainHistory []string
	Pointer      int
}

type CmdStreamFile struct {
	CmdSeeker
	Scanner *bufio.Scanner
}

func (*CmdStreamFile) GetHistory() History {
	return &_NulHistory{}
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

func (stream *CmdStreamFile) ReadLine(ctx context.Context) (context.Context, string, error) {
	if stream.Pointer >= 0 {
		if stream.Pointer < len(stream.PlainHistory) {
			stream.Pointer++
			return ctx, stream.PlainHistory[stream.Pointer-1], nil
		}
		stream.Pointer = -1
	}
	if !stream.Scanner.Scan() {
		if err := stream.Scanner.Err(); err != nil {
			return ctx, "", err
		}
		return ctx, "", io.EOF
	}
	text := strings.TrimRight(stream.Scanner.Text(), "\r\n")
	stream.PlainHistory = append(stream.PlainHistory, text)
	return ctx, text, nil
}

func (sh *Shell) Source(ctx context.Context, fname string) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}
	stream1 := NewCmdStreamFile(fd)
	_, err = sh.Loop(ctx, stream1)
	fd.Close()
	if err == io.EOF {
		return nil
	}
	return err
}
