package conio

import (
	"bytes"
	"strings"
)

func SplitQ(s string) []string {
	args := []string{}
	reader := strings.NewReader(s)
	for reader.Len() > 0 {
		var buffer bytes.Buffer
		for {
			if reader.Len() <= 0 {
				return args
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
		if buffer.Len() > 0 {
			args = append(args, buffer.String())
		}
	}
	return args
}
