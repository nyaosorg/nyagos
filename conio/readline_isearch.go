package conio

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"
)

func KeyFuncIncSearch(this *Buffer) Result {
	var searchBuf bytes.Buffer
	foundStr := ""
	searchStr := ""
	lastDrawWidth := 0
	lastFoundPos := len(this.Session.Histories) - 1
	Backspace(this.Cursor - this.ViewStart)
	for {
		drawStr := fmt.Sprintf("(i-search)[%s]:%s", searchStr, foundStr)
		drawWidth := GetStringWidth(drawStr)
		fmt.Print(drawStr)
		if lastDrawWidth > drawWidth {
			n := lastDrawWidth - drawWidth
			PutRep(' ', n)
			Backspace(n)
		}
		lastDrawWidth = drawWidth
		charcode, _, _ := GetKey()
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
		case '\r': // ENTER
			this.ViewStart = 0
			this.Length = 0
			this.Cursor = 0
			this.ReplaceAndRepaint(0, foundStr)
			return CONTINUE
		case rune('g' & 0x1F):
			w := 0
			var i int
			for i = this.ViewStart; i < this.Cursor; i++ {
				w += GetCharWidth(this.Buffer[i])
				PutRep(this.Buffer[i], 1)
			}
			bs := 0
			for {
				if i >= this.Length {
					if drawWidth > w {
						PutRep(' ', drawWidth-w)
						bs += (drawWidth - w)
					}
					break
				}
				w1 := GetCharWidth(this.Buffer[i])
				if w+w1 >= this.ViewWidth {
					break
				}
				PutRep(this.Buffer[i], 1)
				w += w1
				bs += w1
				i++
			}
			Backspace(bs)
			return CONTINUE
		case rune('r' & 0x1F):
			for i := lastFoundPos - 1; ; i-- {
				if i < 0 {
					i = len(this.Session.Histories) - 1
				}
				if i == lastFoundPos {
					break
				}
				p := this.Session.Histories[i]
				if strings.Contains(p.Line, searchStr) && foundStr != p.Line {
					foundStr = p.Line
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
			for i := len(this.Session.Histories) - 1; ; i-- {
				if i < 0 {
					foundStr = ""
					break
				}
				p := this.Session.Histories[i]
				if strings.Contains(p.Line, searchStr) {
					foundStr = p.Line
					lastFoundPos = i
					break
				}
			}
		}
	}
}
