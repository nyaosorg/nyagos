package nodos

import (
	"os"
	"path/filepath"
	"strings"
)

func JoinList(values ...string) string {
	hash := make(map[string]struct{})

	var buffer strings.Builder
	for _, value := range values {
		for _, val1 := range filepath.SplitList(value) {
			val1 = strings.TrimSpace(val1)
			if len(val1) <= 0 {
				continue
			}
			VAL1 := strings.ToUpper(val1)
			if _, ok := hash[VAL1]; ok {
				continue
			}
			hash[VAL1] = struct{}{}
			if fd, err := os.Open(val1); err != nil {
				continue
			} else {
				fd.Close()
			}
			if buffer.Len() > 0 {
				buffer.WriteRune(os.PathListSeparator)
			}
			buffer.WriteString(val1)
		}
	}
	return buffer.String()
}
