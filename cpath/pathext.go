package cpath

import (
	"os"
	"strings"
)

func IsExecutableSuffix(path string) bool {
	pathExt := os.Getenv("PATHEXT")
	if pathExt != "" {
		for _, ext := range strings.Split(pathExt, ";") {
			if strings.EqualFold(ext, path) {
				return true
			}
		}
	}
	return false
}
