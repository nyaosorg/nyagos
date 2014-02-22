package box

import "bytes"
import "io"
import "strings"

import "github.com/mattn/go-runewidth"

func Print(nodes []string, width int, out io.Writer) {
	maxLen := 1
	for _, finfo := range nodes {
		length := runewidth.StringWidth(finfo)
		if length > maxLen {
			maxLen = length
		}
	}
	nodePerLine := (width - 1) / (maxLen + 1)
	if nodePerLine <= 0 {
		nodePerLine = 1
	}
	nlines := (len(nodes) + nodePerLine - 1) / nodePerLine

	lines := make([]bytes.Buffer, nlines)
	for i, finfo := range nodes {
		lines[i%nlines].WriteString(finfo)
		lines[i%nlines].WriteString(
			strings.Repeat(" ", maxLen+1-
				runewidth.StringWidth(finfo)))
	}
	for _, line := range lines {
		io.WriteString(out, line.String())
		io.WriteString(out, "\n")
	}
}
