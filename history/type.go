package history

import (
	"fmt"
	"os"
	"time"
)

type Line struct {
	Text  string
	Dir   string
	Stamp time.Time
	Pid   int
}

type Container struct {
	rows []Line
}

var NoInstance = &Container{}

func (this *Container) Len() int {
	return len(this.rows)
}

func (this *Container) At(n int) string {
	for n < 0 {
		n += len(this.rows)
	}
	return this.rows[n%len(this.rows)].Text
}

func (this *Container) Push(line string) {
	this.rows = append(this.rows, Line{Text: line})
}

func (this *Container) PushLine(row Line) {
	this.rows = append(this.rows, row)
}

func (row *Line) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%d",
		row.Text,
		row.Dir,
		row.Stamp.Format("2006-01-02 15:04:05"),
		row.Pid)
}

func NewHistoryLine(text string) Line {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}
	return Line{Text: text, Dir: wd, Stamp: time.Now(), Pid: os.Getpid()}
}
