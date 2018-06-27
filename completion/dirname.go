package completion

import "strings"

// DirName returns a directory-part from path, but do not clean up path.
func DirName(path string) string {
	lastroot := strings.LastIndexAny(path, `\/:`)
	if lastroot >= 0 {
		return path[0:(lastroot + 1)]
	}
	return ""
}
