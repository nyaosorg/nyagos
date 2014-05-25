package prompt

import "unicode"
import "os"
import "time"
import "bytes"
import "fmt"

func Format2Prompt(format string) string {
	if format == "" {
		format = "[$P]$_$$$S"
	}
	var buffer bytes.Buffer
	lastchar := '\000'
	for _, ch := range format {
		if lastchar == '$' {
			c := unicode.ToLower(ch)
			if c == 'a' {
				buffer.WriteRune('&')
			} else if c == 'b' {
				buffer.WriteRune('|')
			} else if c == 'c' {
				buffer.WriteRune('(')
			} else if c == 'd' {
				year, month, day := time.Now().Date()
				buffer.WriteString(
					fmt.Sprintf("%04d-%02d-%02d",
						year, month+1, day))
			} else if c == 'e' {
				buffer.WriteRune('\x1B')
			} else if c == 'f' {
				buffer.WriteRune(')')
			} else if c == 'g' {
				buffer.WriteRune('>')
			} else if c == 'h' {
				buffer.WriteRune('\b')
			} else if c == 'l' {
				buffer.WriteRune('<')
			} else if c == 'n' {
				wd, err := os.Getwd()
				if err == nil {
					buffer.WriteString(wd[:2])
				}
			} else if c == 'p' {
				wd, err := os.Getwd()
				if err == nil {
					buffer.WriteString(wd)
				}
			} else if c == 'q' {
				buffer.WriteRune('=')
			} else if c == 's' {
				buffer.WriteRune(' ')
			} else if c == 't' {
				now := time.Now()
				hour, min, sec := now.Clock()
				nnn := now.Nanosecond() / 10000000
				buffer.WriteString(
					fmt.Sprintf("%02d:%02d:%02d.%02d",
						hour, min, sec, nnn))
			} else if c == 'v' {
				// Windows Version
			} else if c == '_' {
				buffer.WriteRune('\n')
			} else if c == '$' {
				buffer.WriteRune('$')
				ch = '\000'
			} else {
				buffer.WriteRune('$')
				buffer.WriteRune(ch)
			}
		} else if ch != '$' {
			buffer.WriteRune(ch)
		}
		lastchar = ch
	}
	return buffer.String()
}
