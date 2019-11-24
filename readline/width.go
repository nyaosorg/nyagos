package readline

import "github.com/zetamatta/go-termgap/hybrid"

type width_t int

var widthCache = map[rune]width_t{}

func ResetCharWidth() {
	widthCache = map[rune]width_t{}
}

func SetCharWidth(c rune, width int) {
	widthCache[c] = width_t(width)
}

func lenEscaped(c rune) width_t {
	w := width_t(3)
	for c > 0xF {
		c >>= 4
		w++
	}
	return w
}

func GetCharWidth(n rune) width_t {
	if n < ' ' {
		return 2
	}
	width, ok := widthCache[n]
	if !ok {
		if n > 0x10000 && !SurrogatePairOk {
			width = lenEscaped(n)
		} else {
			width = width_t(hybrid.RuneWidth(n))
			if width == 0 {
				width = lenEscaped(n)
			}
		}
		widthCache[n] = width
	}
	return width
}

func GetStringWidth(s string) width_t {
	width := width_t(0)
	for _, ch := range s {
		width += GetCharWidth(ch)
	}
	return width
}
