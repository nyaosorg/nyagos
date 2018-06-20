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
	line = rxDoubleQuoted.ReplaceAllStringFunc(line, func(str string) string {
		str = rxSpace.ReplaceAllString(str, "\001")
		str = strings.Replace(str, `'`, "\001", -1)
		return str
	})
	line = rxSingleQuoted.ReplaceAllStringFunc(line, func(str string) string {
		return rxSpace.ReplaceAllString(str, "\001")
	})
	return rxNonSpaces.FindAllStringIndex(line, -1)
}
