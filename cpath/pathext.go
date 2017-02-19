package cpath

import (
	"os"
	"path/filepath"
	"strings"
)

// returns true if suffix exists in %PATHEXT%
func IsExecutableSuffix(path string) bool {
	pathExt := os.Getenv("PATHEXT")
	if pathExt != "" {
		for _, ext := range filepath.SplitList(pathExt) {
			if strings.EqualFold(ext, path) {
				return true
			}
		}
	}
	return false
}
