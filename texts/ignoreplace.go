package texts

import (
	"strings"
)

func ReplaceIgnoreCase(s, old, new string) string {
	var buffer strings.Builder
	for {
		pos := strings.Index(strings.ToUpper(s), strings.ToUpper(old))
		if pos < 0 {
			buffer.WriteString(s)
			return buffer.String()
		}
		buffer.WriteString(s[:pos])
		buffer.WriteString(new)
		s = s[pos+len(old):]
	}
}
