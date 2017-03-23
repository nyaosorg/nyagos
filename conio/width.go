package conio

import "github.com/mattn/go-runewidth"

var widthCache = map[rune]int{}

func ResetCharWidth() {
	widthCache = map[rune]int{}
}

func SetCharWidth(c rune, width int) {
	widthCache[c] = width
}

func GetCharWidth(n rune) int {
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
