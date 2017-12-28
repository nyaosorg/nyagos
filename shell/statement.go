package shell

import (
	"bytes"
	"unicode"
)

func SplitToStatement(line string) []string {
	result := make([]string, 0)
	quote := false
	var buffer bytes.Buffer
	lastc := ' '

	done := func() {
		result = append(result, buffer.String())
		// println("[" + buffer.String() + "]")
		buffer.Reset()
	}

	for _, c := range line {
		if c == '"' {
			quote = !quote
		} else if !quote && c == ';' && unicode.IsSpace(lastc) {
			done()
			lastc = c
			continue
		}
		buffer.WriteRune(c)
		lastc = c
	}
	done()
	return result
}
