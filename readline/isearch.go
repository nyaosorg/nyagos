package readline

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	"github.com/zetamatta/go-getch"
)

func KeyFuncIncSearch(ctx context.Context, this *Buffer) Result {
	var searchBuf strings.Builder
	foundStr := ""
	searchStr := ""
	lastDrawWidth := 0
	lastFoundPos := this.History.Len() - 1
	Backspace(this.Cursor - this.ViewStart)

	update := func() {
		for i := this.History.Len() - 1; ; i-- {
			if i < 0 {
				foundStr = ""
				break
			}
			line := this.History.At(i)
			if strings.Contains(line, searchStr) {
				foundStr = line
				lastFoundPos = i
				break
			}
		}
	}
	for {
		drawStr := fmt.Sprintf("(i-search)[%s]:%s", searchStr, foundStr)
		drawWidth := 0
		for _, ch := range drawStr {
			w1 := GetCharWidth(ch)
			if drawWidth+w1 >= this.ViewWidth() {
				break
			}
			PutRune(ch)
			drawWidth += w1
		}
		if lastDrawWidth > drawWidth {
			n := lastDrawWidth - drawWidth
			PutRunes(' ', n)
			Backspace(n)
		}
		lastDrawWidth = drawWidth
		fmt.Fprint(Console, CURSOR_ON)
		Console.Flush()
		charcode := getch.Rune()
		fmt.Fprint(Console, CURSOR_OFF)
		Backspace(drawWidth)
		switch charcode {
		case '\b':
			searchBuf.Reset()
			// chop last char
			var lastchar rune
			for i, c := range searchStr {
				if i > 0 {
					searchBuf.WriteRune(lastchar)
				}
				lastchar = c
			}
			searchStr = searchBuf.String()
			update()
		case '\r': // ENTER
			this.ViewStart = 0
			this.Length = 0
			this.Cursor = 0
			this.ReplaceAndRepaint(0, foundStr)
			return CONTINUE
		case rune('c' & 0x1F), rune('g' & 0x1F), rune(0x1B):
			w := 0
			var i int
			for i = this.ViewStart; i < this.Cursor; i++ {
				w += GetCharWidth(this.Buffer[i])
				PutRune(this.Buffer[i])
			}
			bs := 0
			for {
				if i >= this.Length {
					if drawWidth > w {
						PutRunes(' ', drawWidth-w)
						bs += (drawWidth - w)
					}
					break
				}
				w1 := GetCharWidth(this.Buffer[i])
				if w+w1 >= this.ViewWidth() {
					break
				}
				PutRune(this.Buffer[i])
				w += w1
				bs += w1
				i++
			}
			Backspace(bs)
			return CONTINUE
		case rune('r' & 0x1F):
			for i := lastFoundPos - 1; ; i-- {
				if i < 0 {
					i = this.History.Len() - 1
				}
				if i == lastFoundPos {
					break
				}
				line := this.History.At(i)
				if strings.Contains(line, searchStr) && foundStr != line {
					foundStr = line
					lastFoundPos = i
					break
				}
			}
		case rune('s' & 0x1F):
			for i := lastFoundPos + 1; ; i++ {
				if i >= this.History.Len() {
					i = 0
					break
				}
				if i == lastFoundPos {
					break
				}
				line := this.History.At(i)
				if strings.Contains(line, searchStr) && foundStr != line {
					foundStr = line
					lastFoundPos = i
					break
				}
			}
		default:
			if unicode.IsControl(charcode) {
				break
			}
			searchBuf.WriteRune(charcode)
			searchStr = searchBuf.String()
			update()
		}
	}
}
