package dos

import "strings"

func YenYen2Yen(path string) string {
	result := strings.Replace(path, "\\\\", "\\", -1)
	if strings.HasPrefix(path, "\\\\") {
		return "\\" + result
	} else {
		return result
	}
}
