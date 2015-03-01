package conio

import (
	"bytes"
	"strings"
)

func QuotedWordCutter(reader *strings.Reader) (string, bool) {
	var buffer bytes.Buffer
	for {
		if reader.Len() <= 0 {
			return "", false
		}
		ch, _, _ := reader.ReadRune()
		if ch != ' ' {
			reader.UnreadRune()
			break
		}
	}
	quote := false
	for reader.Len() > 0 {
		ch, _, _ := reader.ReadRune()
		if ch == '"' {
			quote = !quote
		}
		if ch == ' ' && !quote {
			break
		}
		buffer.WriteRune(ch)
	}
	return buffer.String(), true
}

// Split s with SPACES not enclosing with double-quotations.
func SplitQ(line string) []string {
	args := make([]string, 0, 10)
	reader := strings.NewReader(line)
	for reader.Len() > 0 {
		word, ok := QuotedWordCutter(reader)
		if ok {
			args = append(args, word)
		}
	}
	return args
}

func QuotedFirstWord(line string) string {
	reader := strings.NewReader(line)
	str, _ := QuotedWordCutter(reader)
	return str
}
