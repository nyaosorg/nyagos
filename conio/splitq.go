package conio

import (
	"bytes"
	"strings"
)

const NULQUOTE = '\000'

func quotedWordCutter(reader *strings.Reader) (string, bool) {
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
	quote := NULQUOTE
	yenCount := 0
	for reader.Len() > 0 {
		ch, _, _ := reader.ReadRune()
		if yenCount%2 == 0 {
			if quote == NULQUOTE && (ch == '"' || ch == '\'') {
				quote = ch
			} else if quote != NULQUOTE && ch == quote {
				quote = NULQUOTE
			}
		}
		if ch == ' ' && quote == NULQUOTE {
			break
		}
		if ch == '\\' {
			yenCount++
		} else {
			yenCount = 0
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
		word, ok := quotedWordCutter(reader)
		if ok {
			args = append(args, word)
		}
	}
	return args
}

func QuotedFirstWord(line string) string {
	reader := strings.NewReader(line)
	str, _ := quotedWordCutter(reader)
	return str
}
