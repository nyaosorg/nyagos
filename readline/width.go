package readline

import "github.com/mattn/go-runewidth"

var widthCache = map[rune]int{}

func ResetCharWidth() {
	widthCache = map[rune]int{}
}

func SetCharWidth(c rune, width int) {
	widthCache[c] = width
}

func lenEscaped(c rune) int {
	w := 3
	for c > 0xF {
		c >>= 4
		w++
	}
	return w
}

func GetCharWidth(n rune) int {
	if n < ' ' {
		return 2
	}
	width, ok := widthCache[n]
	if !ok {
		if n > 0x10000 && !SurrogatePairOk {
			width = lenEscaped(n)
		} else {
			width = runewidth.RuneWidth(n)
			if width == 0 {
				width = lenEscaped(n)
			}
		}
		widthCache[n] = width
	}
	return width
}

func GetStringWidth(s string) int {
	width := 0
	for _, ch := range s {
		width += GetCharWidth(ch)
	}
	return width
}
