package frame

import (
	"unicode"

	"github.com/nyaosorg/go-readline-ny"
)

type _Coloring struct {
	bits        int
	last        rune
	defaultBits int
}

func (s *_Coloring) Init() int {
	s.bits = s.defaultBits
	return readline.White
}

const (
	backquotedBit = 1
	percentBit    = 2
	quotedBit     = 4
	optionBit     = 8
)

func (s *_Coloring) Next(codepoint rune) int {
	newbits := s.bits
	if codepoint == '`' {
		newbits ^= backquotedBit
	} else if codepoint == '%' {
		newbits ^= percentBit
	} else if codepoint == '"' {
		newbits ^= quotedBit
	} else if s.last == ' ' && (codepoint == '/' || codepoint == '-') {
		newbits ^= optionBit
	} else if (s.bits&optionBit) != 0 && !unicode.IsLetter(codepoint) && (codepoint < '0' || codepoint > '9') && codepoint != '-' {
		newbits &^= optionBit
	} else if s.last == '%' && (s.bits&percentBit) != 0 && unicode.IsDigit(codepoint) {
		newbits &^= percentBit
	}
	bits := s.bits | newbits
	color := readline.White
	if unicode.IsControl(codepoint) {
		color = readline.Blue
	} else if codepoint == '\u3000' {
		color = 41
	} else if (bits & percentBit) != 0 {
		color = readline.Cyan
	} else if (bits & backquotedBit) != 0 {
		color = readline.Red
	} else if (bits & quotedBit) != 0 {
		color = readline.Magenta
	} else if (newbits & optionBit) != 0 {
		color = readline.DarkYellow
	} else if codepoint == '&' || codepoint == '|' || codepoint == '<' || codepoint == '>' || (s.last == ' ' && codepoint == ';') {
		color = readline.Green
	} else {
		color = readline.White
	}
	s.bits = newbits
	s.last = codepoint
	return color
}
