package texts

import (
	"regexp"
	"strings"
)

var rxEscape = regexp.MustCompile(`\\["'\\]`)
var rxDoubleQuoted = regexp.MustCompile(`"[^"]*"`)
var rxSingleQuoted = regexp.MustCompile(`'[^']*'`)
var rxSpace = regexp.MustCompile(`\s`)
var rxNonSpaces = regexp.MustCompile(`\S+`)

// SplitLikeShell splits string with spaces not enclused with double-quotations. This version returns fields' indexes rather than themselves.
func SplitLikeShell(line string) [][]int {
	line = rxEscape.ReplaceAllString(line, "\001\001")

	var buffer strings.Builder
	quoted := false
	for _, r := range line {
		if r == '"' {
			quoted = !quoted
		}
		if quoted && (r == ' ' || r == '\'') {
			r = '\001'
		}
		buffer.WriteRune(r)
	}
	line = buffer.String()

	line = rxSingleQuoted.ReplaceAllStringFunc(line, func(str string) string {
		return rxSpace.ReplaceAllString(str, "\001")
	})
	return rxNonSpaces.FindAllStringIndex(line, -1)
}
