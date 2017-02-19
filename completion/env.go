package completion

import (
	"os"
	"strings"

	"../interpreter"
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

func listUpOsEnv(name string) []string {
	matches := []string{}
	for _, envEquation := range os.Environ() {
		equalPos := strings.IndexRune(envEquation, '=')
		if equalPos >= 0 {
			envName := envEquation[:equalPos]
			if strings.HasPrefix(strings.ToUpper(envName), name) {
				matches = append(matches, "%"+envName+"%")
			}
		}
	}
	return matches
}

func listUpDynamicEnv(name string) []string {
	matches := []string{}
	for envName, _ := range interpreter.PercentFunc {
		if strings.HasPrefix(envName, name) {
			matches = append(matches, "%"+envName+"%")
		}
	}
	return matches
}

var PercentFuncs = []func(string) []string{
	listUpOsEnv,
	listUpDynamicEnv,
}

func listUpEnv(cmdline string) ([]string, int, error) {
	lastPercentPos := findLastOpenPercent(cmdline)
	if lastPercentPos < 0 {
		return nil, -1, nil
	}
	str := cmdline[lastPercentPos:]
	name := strings.ToUpper(str[1:])
	var matches []string
	for _, f := range PercentFuncs {
		matches = append(matches, f(name)...)
	}
	if len(matches) <= 0 { // nothing matches.
		return nil, -1, nil
	}
	return matches, lastPercentPos, nil
}
