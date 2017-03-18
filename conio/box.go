package conio

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
	getch "github.com/zetamatta/go-getch"
)

var ansiCutter = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")

func BoxPrint(ctx context.Context, nodes []string, out io.Writer) bool {
	value, _, _ := boxPrint(ctx, nodes, 0, false, out)
	return value
}

func boxPrint(ctx context.Context, nodes []string, offset int, paging bool, out io.Writer) (bool, int, int) {
	csbi := GetScreenBufferInfo()
	width := int(csbi.Size.X)
	if width <= 0 || width > 999 {
		width = 80
	}
	maxLen := 1
	for _, finfo := range nodes {
		length := runewidth.StringWidth(ansiCutter.ReplaceAllString(finfo, ""))
		if length > maxLen {
			maxLen = length
		}
	}
	nodePerLine := (width - 1) / (maxLen + 1)
	if nodePerLine <= 0 {
		nodePerLine = 1
	}
	nlines := (len(nodes) + nodePerLine - 1) / nodePerLine

	lines := make([][]byte, nlines)
	row := 0
	for _, finfo := range nodes {
		lines[row] = append(lines[row], finfo...)
		w := runewidth.StringWidth(ansiCutter.ReplaceAllString(finfo, ""))
		for i, iEnd := 0, maxLen+1-w; i < iEnd; i++ {
			lines[row] = append(lines[row], ' ')
		}
		row++
		if row >= nlines {
			row = 0
		}
	}
	i_end := len(lines)
	if paging {
		_, height := csbi.ViewSize()
		height--
		if i_end >= offset+height {
			i_end = offset + height
		}
	}

	for i := offset; i < i_end; i++ {
		fmt.Fprintln(out, string(lines[i]))
		if ctx != nil {
			select {
			case <-ctx.Done():
				return false, nodePerLine, nlines
			default:
			}
		}
	}
	return true, nodePerLine, nlines
}

const (
	CURSOR_OFF = "\x1B[?25l"
	CURSOR_ON  = "\x1B[?25h"
	BOLD_ON    = "\x1B[0;47;30m"
	BOLD_OFF   = "\x1B[0m"

	K_LEFT  = 0x25
	K_RIGHT = 0x27
	K_UP    = 0x26
	K_DOWN  = 0x28
)

func truncate(s string, w int) string {
	return runewidth.Truncate(strings.TrimSpace(s), w, "")
}

func BoxChoice(nodes []string, out io.Writer) string {
	cursor := 0
	nodes_draw := make([]string, len(nodes))
	width, height := GetScreenBufferInfo().ViewSize()
	width--
	height--
	for i := 0; i < len(nodes); i++ {
		nodes_draw[i] = truncate(nodes[i], width-1)
	}
	io.WriteString(out, CURSOR_OFF)
	defer io.WriteString(out, CURSOR_ON)

	offset := 0
	for {
		nodes_draw[cursor] = BOLD_ON +
			truncate(nodes[cursor], width-1) + BOLD_OFF
		status, _, h := boxPrint(nil, nodes_draw, offset, true, out)
		if !status {
			return ""
		}
		nodes_draw[cursor] = truncate(nodes[cursor], width-1)
		last := cursor
		for last == cursor {
			e := getch.All()
			if k := e.Key; k != nil {
				switch k.Rune {
				case 'h', ('b' & 0x1F):
					if cursor-h >= 0 {
						cursor -= h
					}
				case 'l', ('f' & 0x1F):
					if cursor+h < len(nodes) {
						cursor += h
					}
				case 'j', ('n' & 0x1F), ' ':
					if cursor+1 < len(nodes) {
						cursor++
					}
				case 'k', ('p' & 0x1F), '\b':
					if cursor > 0 {
						cursor--
					}
				case '\r', '\n':
					return nodes[cursor]
				case '\x1B', ('g' & 0x1F):
					return ""
				}

				switch k.Scan {
				case K_LEFT:
					if cursor-h >= 0 {
						cursor -= h
					}
				case K_RIGHT:
					if cursor+h < len(nodes) {
						cursor += h
					}
				case K_DOWN:
					if cursor+1 < len(nodes) {
						cursor++
					}
				case K_UP:
					if cursor > 0 {
						cursor--
					}
				}

				// x := cursor / h
				y := cursor % h
				if y < offset {
					offset--
				} else if y >= offset+height {
					offset++
				}
			}
		}
		if h < height {
			fmt.Fprintf(out, "\x1B[%dA", h)
		} else {
			fmt.Fprintf(out, "\x1B[%dA", height)
		}
	}
}
