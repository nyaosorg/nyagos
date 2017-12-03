package shell

import (
	"context"
	"io"
)

type BufStream struct {
	line []string
	n    int
}

func (this *BufStream) ReadLine(c context.Context) (context.Context, string, error) {
	if this.n >= len(this.line) {
		return c, "", io.EOF
	}
	this.n++
	return c, this.line[this.n-1], nil
}

func (this *BufStream) GetPos() int {
	return this.n
}

func (this *BufStream) SetPos(n int) error {
	this.n = n
	return nil
}

func (this *BufStream) Add(line string) {
	this.line = append(this.line, line)
}

func (this *Cmd) Eval(c context.Context, line string) (errorlevel int, err error) {
	bufStream := BufStream{}
	bufStream.Add(line)

	rc, err := this.Loop(&bufStream)
	if err == nil || err == io.EOF {
		return rc, nil
	} else {
		return rc, err
	}
}
