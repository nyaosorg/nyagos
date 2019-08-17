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
			if len(val1) > 0 {
				VAL1 := strings.ToUpper(val1)
				if _, ok := hash[VAL1]; !ok {
					hash[VAL1] = struct{}{}
					if buffer.Len() > 0 {
						buffer.WriteRune(os.PathListSeparator)
					}
					buffer.WriteString(val1)
				}
			}
		}
	}
	return buffer.String()
}
