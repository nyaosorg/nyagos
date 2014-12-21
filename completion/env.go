package completion

import (
	"os"
	"strings"
)

func findLastOpenPercent(this string) int {
	pos := -1
	for i, ch := range this {
		if ch == '%' {
			if pos >= 0 { // closing percent mark
				pos = -1
			} else { // opening percent mark
				pos = i
			}
		}
	}
	return pos
}

func listUpEnv(cmdline string) ([]string, int, error) {
	lastPercentPos := findLastOpenPercent(cmdline)
	if lastPercentPos < 0 {
		return nil, -1, nil
	}
	str := cmdline[lastPercentPos:]
	name := strings.ToUpper(str[1:])
	matches := make([]string, 0, 5)
	for _, envEquation := range os.Environ() {
		equalPos := strings.IndexRune(envEquation, '=')
		if equalPos >= 0 {
			envName := envEquation[:equalPos]
			if strings.HasPrefix(strings.ToUpper(envName), name) {
				matches = append(matches, "%"+envName+"%")
			}
		}
	}
	if len(matches) <= 0 { // nothing matches.
		return nil, -1, nil
	}
	return matches, lastPercentPos, nil
}
