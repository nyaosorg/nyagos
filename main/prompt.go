package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"../cpath"
	"../dos"
)

func Format2Prompt(format string) string {
	if format == "" {
		format = "[$P]$_$$$S"
	}
	var buffer bytes.Buffer
	lastchar := '\000'
	for reader := strings.NewReader(format); reader.Len() > 0; {
		ch, _, err := reader.ReadRune()
		if err != nil {
			break
		}
		if lastchar == '$' {
			c := unicode.ToLower(ch)
			if c == 'a' {
				buffer.WriteRune('&')
			} else if c == 'b' {
				buffer.WriteRune('|')
			} else if c == 'c' {
				buffer.WriteRune('(')
			} else if c == 'd' {
				buffer.WriteString(time.Now().Format("2006-01-02"))
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
				if wd, err := dos.Getwd(); err != nil {
					fmt.Fprintf(os.Stderr, "$P: %s\n", err.Error())
				} else {
					buffer.WriteString(cpath.ReplaceHomeToTildeSlash(wd))
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
			} else if c == 'u' {
				r := 0
				for i := 0; i < 4 && reader.Len() > 0; i++ {
					r1, _, err := reader.ReadRune()
					if err != nil {
						break
					}
					n := strings.IndexRune("0123456789ABCDEF",
						unicode.ToUpper(r1))
					if n < 0 {
						reader.UnreadRune()
						break
					}
					r = r*16 + n
				}
				if r > 0 {
					buffer.WriteRune(rune(r))
				}
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
