package findfile

import (
	"path/filepath"
	"strings"
)

// Expand filenames matching with wildcard-pattern.
func Glob(pattern string) ([]string, error) {
	pname := filepath.Base(pattern)
	if strings.IndexAny(pname, "*?") < 0 {
		return nil, nil
	}
	match := make([]string, 0, 100)
	dirname := filepath.Dir(pattern)
	err := Walk(pattern, func(findf *FileInfo) bool {
		name := findf.Name()
		if (name[0] != '.' || pname[0] == '.') && !findf.IsHidden() {
			match = append(match, filepath.Join(dirname, name))
		}
		return true
	})
	return match, err
}

func Globs(patterns []string) []string {
	result := make([]string, 0, len(patterns))
	for _, pattern1 := range patterns {
		matches, err := Glob(pattern1)
		if matches == nil || len(matches) <= 0 || err != nil {
			result = append(result, pattern1)
		} else {
			for _, s := range matches {
				result = append(result, s)
			}
		}
	}
	return result
}
