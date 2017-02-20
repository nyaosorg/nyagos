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

func listUpOsEnv(name string) []Element {
	matches := []Element{}
	for _, envEquation := range os.Environ() {
		equalPos := strings.IndexRune(envEquation, '=')
		if equalPos >= 0 {
			envName := envEquation[:equalPos]
			if strings.HasPrefix(strings.ToUpper(envName), name) {
				envValue := "%" + envName + "%"
				element := Element{InsertStr: envValue, ListupStr: envValue}
				matches = append(matches, element)
			}
		}
	}
	return matches
}

func listUpEnv(cmdline string) ([]Element, int, error) {
	lastPercentPos := findLastOpenPercent(cmdline)
	if lastPercentPos < 0 {
		return nil, -1, nil
	}
	str := cmdline[lastPercentPos:]
	name := strings.ToUpper(str[1:])
	var matches []Element
	for _, f := range PercentFuncs {
		matches = append(matches, f(name)...)
	}
	if len(matches) <= 0 { // nothing matches.
		return nil, -1, nil
	}
	return matches, lastPercentPos, nil
}
