package history

type Row struct {
	Text string
	Dir  string
}

type THistory struct {
	rows []Row
}

func (this *THistory) Len() int {
	return len(this.rows)
}

func (this *THistory) At(n int) string {
	for n < 0 {
		n += len(this.rows)
	}
	return this.rows[n%len(this.rows)].Text
}

func (this *THistory) Push(line string) {
	this.rows = append(this.rows, Row{Text: line})
}

func (this *THistory) PushRow(row Row) {
	this.rows = append(this.rows, row)
}
