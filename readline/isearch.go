package readline

import (
	"context"
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"
)

func KeyFuncIncSearch(ctx context.Context, this *Buffer) Result {
	var searchBuf strings.Builder
	foundStr := ""
	searchStr := ""
	lastFoundPos := this.History.Len() - 1
	this.Backspace(this.Cursor - this.ViewStart)

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
			this.PutRune(ch)
			drawWidth += w1
		}
		this.Eraseline()
		io.WriteString(this.Out, CURSOR_ON)
		this.Out.Flush()
		key, err := getKey(this.TTY)
		if err != nil {
			println(err.Error())
			return CONTINUE
		}
		io.WriteString(this.Out, CURSOR_OFF)
		this.Backspace(drawWidth)
		switch key {
		case "\b":
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
		case "\r":
			this.ViewStart = 0
			this.Length = 0
			this.Cursor = 0
			this.ReplaceAndRepaint(0, foundStr)
			return CONTINUE
		case "\x03", "\x07", "\x1B":
			w := 0
			var i int
			for i = this.ViewStart; i < this.Cursor; i++ {
				w += GetCharWidth(this.Buffer[i])
				this.PutRune(this.Buffer[i])
			}
			bs := 0
			for {
				if i >= this.Length {
					if drawWidth > w {
						this.PutRunes(' ', drawWidth-w)
						bs += (drawWidth - w)
					}
					break
				}
				w1 := GetCharWidth(this.Buffer[i])
				if w+w1 >= this.ViewWidth() {
					break
				}
				this.PutRune(this.Buffer[i])
				w += w1
				bs += w1
				i++
			}
			this.Backspace(bs)
			return CONTINUE
		case "\x12":
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
		case "\x13":
			for i := lastFoundPos + 1; ; i++ {
				if i >= this.History.Len() {
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
			charcode, _ := utf8.DecodeRuneInString(key)
			if unicode.IsControl(charcode) {
				break
			}
			searchBuf.WriteRune(charcode)
			searchStr = searchBuf.String()
			update()
		}
	}
}
