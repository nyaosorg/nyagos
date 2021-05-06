package shell

import (
	"strings"
)

func Quote(src string) string {
	var buffer strings.Builder
	buffer.WriteByte('"')
	backSlashCount := 0
	for _, c := range src {
		if c == '\\' {
			backSlashCount++
			continue
		}
		if c == '"' {
			for ; backSlashCount > 0; backSlashCount-- {
				buffer.Write([]byte{'\\', '\\'})
			}
			buffer.Write([]byte{'\\', '"'})
			continue
		}
		for ; backSlashCount > 0; backSlashCount-- {
			buffer.WriteByte('\\')
		}
		buffer.WriteRune(c)
	}
	for ; backSlashCount > 0; backSlashCount-- {
		buffer.Write([]byte{'\\', '\\'})
	}
	buffer.WriteByte('"')
	return buffer.String()
}
