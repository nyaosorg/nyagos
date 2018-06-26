package history

import (
	"fmt"
	"os"
	"time"
)

// Line has one history data
type Line struct {
	Text  string
	Dir   string
	Stamp time.Time
	Pid   int
}

// Container has all history data.
type Container struct {
	rows []Line
}

type packageIdT struct{}

// PackageId is the unique mark to use as Context key
var PackageId packageIdT

// Len returns size of history
func (c *Container) Len() int {
	return len(c.rows)
}

// At returns n-th history-text
func (c *Container) At(n int) string {
	for n < 0 {
		n += len(c.rows)
	}
	return c.rows[n%len(c.rows)].Text
}

// Push appends a new history line to self with string
func (c *Container) Push(line string) {
	c.rows = append(c.rows, Line{Text: line})
}

// PushLine appends a new history line to self with Line object
func (c *Container) PushLine(row Line) {
	c.rows = append(c.rows, row)
}

// String returns self as printable text
func (row *Line) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%d",
		row.Text,
		row.Dir,
		row.Stamp.Format("2006-01-02 15:04:05"),
		row.Pid)
}

// NewHistoryLine returns new Line object with history-text
func NewHistoryLine(text string) Line {
	wd, err := os.Getwd()
	if err != nil {
		wd = ""
	}
	return Line{Text: text, Dir: wd, Stamp: time.Now(), Pid: os.Getpid()}
}
