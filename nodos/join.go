package nodos

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var rxDrive = regexp.MustCompile("^[a-zA-Z]:")

func joinPath2(a, b string) string {
	if len(a) <= 0 || strings.HasPrefix(b, `\\`) || rxDrive.MatchString(b) {
		return b
	}
	if b[0] == '\\' || b[0] == '/' {
		if rxDrive.MatchString(a) {
			return a[:2] + b
		}
		return b
	}
	switch a[len(a)-1] {
	case '\\', '/', ':':
		return a + b
	default:
		return fmt.Sprintf("%s%c%s", a, os.PathSeparator, b)
	}
}

// Join is compatible with CPath::Combine of MFC (ex:`C:\foo` + `\bar` -> `c:\bar`)
// Do not clean path (keep `./` on arguments)
func Join(paths ...string) string {
	result := paths[len(paths)-1]
	for i := len(paths) - 2; i >= 0; i-- {
		result = joinPath2(paths[i], result)
	}
	return result
}
