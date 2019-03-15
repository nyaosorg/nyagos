package readline

import "github.com/mattn/go-runewidth"

var specialRune = map[rune]string{
	0x200D: "<ZWJ>",
}

var widthCache = map[rune]int{}

func ResetCharWidth() {
	widthCache = map[rune]int{}
}

func SetCharWidth(c rune, width int) {
	widthCache[c] = width
}

func GetCharWidth(n rune) int {
	if n < ' ' {
		return 2
	}
	if text, ok := specialRune[n]; ok {
		return len(text)
	}
	width, ok := widthCache[n]
	if !ok {
		width = runewidth.RuneWidth(n)
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
