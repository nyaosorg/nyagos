package frame

import (
	"github.com/nyaosorg/go-readline-ny"
)

type _Coloring struct {
	bits int
}

func (s *_Coloring) Init() {
	s.bits = 0
}

func (s *_Coloring) Get(codepoint rune) int {
	const (
		backquotedBit = 1
		percentBit    = 2
		quotedBit     = 4
	)
	newbits := s.bits
	if codepoint == '`' {
		newbits ^= backquotedBit
	} else if codepoint == '%' {
		newbits ^= percentBit
	} else if codepoint == '"' {
		newbits ^= quotedBit
	}
	bits := s.bits | newbits
	color := readline.White
	if (bits & backquotedBit) != 0 {
		color = readline.Red
	} else if (bits & percentBit) != 0 {
		color = readline.Cyan
	} else if (bits & quotedBit) != 0 {
		color = readline.Magenta
	} else {
		color = readline.White
	}
	s.bits = newbits
	return color
}
