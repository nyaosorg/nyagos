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
	value, _, _ := boxPrint(ctx, nodes, out)
	return value
}

func boxPrint(ctx context.Context, nodes []string, out io.Writer) (bool, int, int) {
	width := int(GetScreenBufferInfo().Size.X)
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
	for _, line := range lines {
		fmt.Fprintln(out, strings.TrimSpace(string(line)))
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

func BoxChoice(nodes []string, out io.Writer) string {
	cursor := 0
	nodes_draw := make([]string, len(nodes))
	for i := 0; i < len(nodes); i++ {
		nodes_draw[i] = nodes[i]
	}
	io.WriteString(out, CURSOR_OFF)
	defer io.WriteString(out, CURSOR_ON)
	for {
		nodes_draw[cursor] = BOLD_ON + nodes[cursor] + BOLD_OFF
		status, _, h := boxPrint(nil, nodes_draw, out)
		if !status {
			return ""
		}
		nodes_draw[cursor] = nodes[cursor]
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
		}
		fmt.Fprintf(out, "\x1B[%dA", h)
	}
}
