package history

import "time"

type Line struct {
	Text  string
	Dir   string
	Stamp time.Time
}

type Container struct {
	rows []Line
}

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
