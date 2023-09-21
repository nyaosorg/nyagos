//go:build orgxwidth

package textwidth

import (
	"unicode"

	"golang.org/x/text/width"
)

func newRuneWidth(ambiguousIsWide bool) func(rune) int {
	var aw int
	if ambiguousIsWide {
		aw = 2
	} else {
		aw = 1
	}
	return func(r rune) int {
		if !unicode.IsPrint(r) {
			return 0
		}
		switch width.LookupRune(r).Kind() {
		case width.Neutral, width.EastAsianNarrow, width.EastAsianHalfwidth:
			return 1
		case width.EastAsianWide, width.EastAsianFullwidth:
			return 2
		case width.EastAsianAmbiguous:
			return aw
		default:
			return 0
		}
	}
}
