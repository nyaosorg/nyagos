//go:build !orgxwidth

package textwidth

import (
	"github.com/mattn/go-runewidth"
)

func newRuneWidth(ambiguousIsWide bool) func(rune) int {
	c := runewidth.NewCondition()
	c.EastAsianWidth = ambiguousIsWide
	return c.RuneWidth
}
