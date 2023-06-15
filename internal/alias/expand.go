package alias

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var paramMatch = regexp.MustCompile(`\$(\~)?(\*|[0-9]+)`)

func ExpandMacro(base string, args []string, rawargs []string) string {
	isReplaced := false
	cmdline := paramMatch.ReplaceAllStringFunc(base, func(s string) string {
		if s == "$~*" {
			isReplaced = true
			if len(args) >= 2 {
				return strings.Join(args[1:], " ")
			}
			return ""
		} else if s == "$*" {
			isReplaced = true
			if len(args) >= 2 {
				return strings.Join(rawargs[1:], " ")
			}
			return ""
		} else if len(s) >= 3 && s[0] == '$' && s[1] == '~' && unicode.IsDigit(rune(s[2])) {
			i, err := strconv.ParseInt(s[2:], 10, 0)
			if err == nil {
				isReplaced = true
				if 0 <= i && args != nil && int(i) < len(args) {
					return args[i]
				}
				return ""
			}
		}
		i, err := strconv.ParseInt(s[1:], 10, 0)
		if err == nil {
			isReplaced = true
			if 0 <= i && args != nil && int(i) < len(args) {
				return rawargs[i]
			}
			return ""
		}
		return s
	})

	if !isReplaced {
		var buffer strings.Builder
		buffer.WriteString(base)
		for _, s := range rawargs[1:] {
			fmt.Fprintf(&buffer, " %s", s)
		}
		cmdline = buffer.String()
	}
	return cmdline
}
