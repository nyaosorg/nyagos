package shell

import (
	"context"
	"io"
)

type BufStream struct {
	line    []string
	n       int
	History History
}

func (bufStream *BufStream) ReadLine(c context.Context) (context.Context, string, error) {
	if bufStream.n >= len(bufStream.line) {
		return c, "", io.EOF
	}
	bufStream.n++
	return c, bufStream.line[bufStream.n-1], nil
}

func (bufStream *BufStream) SetPos(n int) error {
	bufStream.n = n
	return nil
}

func (bufStream *BufStream) Add(line string) {
	bufStream.line = append(bufStream.line, line)
}

func (bufStream *BufStream) GetHistory() History {
	return bufStream.History
}
