package conio

import "github.com/mattn/go-runewidth"

var widthCache = map[rune]int{
	rune(0x262D): 2, // Hammer and sickle: http://unicode-table.com/en/262D/
}

func GetCharWidth(n rune) int {
	width, ok := widthCache[n]
	if !ok {
		width = runewidth.RuneWidth(n)
		widthCache[n] = width
	}
	return width
	// if n > 0xFF {
	//	return 2;
	//}else{
	//	return 1;
	//}
}

func GetStringWidth(s string) int {
	width := 0
	for _, ch := range s {
		width += GetCharWidth(ch)
	}
	return width
}
