package completion

import (
	"os"
	"strings"
)

type IVariable interface {
	Lookup(name string) string
	EachKey(func(string))
}

type EnvironmentVariable struct {
}

func (this *EnvironmentVariable) Lookup(name string) string {
	return os.Getenv(name)
}

func (this *EnvironmentVariable) EachKey(f func(name string)) {
	for _, envEquation := range os.Environ() {
		equalPos := strings.IndexRune(envEquation, '=')
		if equalPos >= 0 {
			envName := envEquation[:equalPos]
			f(envName)
		}
	}
}

var PercentVariables = []IVariable{
	new(EnvironmentVariable),
}

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

func listUpEnv(cmdline string) ([]Element, int, error) {
	lastPercentPos := findLastOpenPercent(cmdline)
	if lastPercentPos < 0 {
		return nil, -1, nil
	}
	str := cmdline[lastPercentPos:]
	name := strings.ToUpper(str[1:])
	var matches []Element

	for _, vars := range PercentVariables {
		vars.EachKey(func(envName string) {
			if strings.HasPrefix(strings.ToUpper(envName), name) {
				envValue := "%" + envName + "%"
				element := Element{InsertStr: envValue, ListupStr: envValue}
				matches = append(matches, element)
			}
		})
	}
	if len(matches) <= 0 { // nothing matches.
		return nil, -1, nil
	}
	return matches, lastPercentPos, nil
}
