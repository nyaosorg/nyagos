package nodos

import (
	"os"
	"path/filepath"
	"strings"
)

var pathExtEnvNames = []string{"PATHEXT", "NYAGOSPATHEXT"}

// IsExecutableSuffix returns true if suffix exists in %PATHEXT% and %NYAGOSPATHEXT%
func IsExecutableSuffix(path string) bool {
	for _, envName := range pathExtEnvNames {
		pathExt := os.Getenv(envName)
		if pathExt != "" {
			for _, ext := range filepath.SplitList(pathExt) {
				if strings.EqualFold(ext, path) {
					return true
				}
			}
		}
	}
	return false
}
