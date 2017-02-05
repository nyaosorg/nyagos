package history

type THistory struct {
	body []string
}

func (this *THistory) Len() int {
	return len(this.body)
}

func (this *THistory) At(n int) string {
	for n < 0 {
		n += len(this.body)
	}
	return this.body[n%len(this.body)]
}

func (this *THistory) Push(line string) {
	this.body = append(this.body, line)
}
