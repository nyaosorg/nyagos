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
	} else {
		return err
	}
}
