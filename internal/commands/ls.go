package commands

import (
	"strings"

	"github.com/dustin/go-humanize"
)

func formatByHumanize(size int64) string {
	s := humanize.Bytes(uint64(size))
	if len(s) > 0 && s[len(s)-1] == 'B' {
		s = s[:len(s)-1]
	}
	return strings.ToUpper(strings.ReplaceAll(s, " ", ""))
}
