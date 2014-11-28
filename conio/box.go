package conio

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mattn/go-runewidth"
)

var ansiCutter = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")

func BoxPrint(nodes []string, out io.Writer) {
	width, _ := GetScreenSize() // ignore height
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
	}
}
