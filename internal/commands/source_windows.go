package commands

import (
	"os"
	"strings"
)

func findBatch(name string) (string, bool) {
	lowerName := strings.ToLower(name)
	if strings.HasSuffix(lowerName, ".cmd") || strings.HasSuffix(lowerName, ".bat") {
		return name, true
	}
	tmp := name + ".cmd"
	if _, err := os.Stat(tmp); err == nil {
		return tmp, true
	}
	tmp = name + ".bat"
	if _, err := os.Stat(tmp); err == nil {
		return tmp, true
	}
	return "", false
}
