package completion

import (
	"os"
	"regexp"
	"strings"
)

type IVariable interface {
	Lookup(name string) string
	EachKey(func(string))
}

type EnvironmentVariable struct {
}

func (*EnvironmentVariable) Lookup(name string) string {
	return os.Getenv(name)
}

func (*EnvironmentVariable) EachKey(f func(name string)) {
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

var rxDollar = regexp.MustCompile(`\$\w+$`)
var rxDollar2 = regexp.MustCompile(`\$\{[^\}]+$`)

func listUpEnv(cmdline string) ([]Element, int, error) {
	var makeCandidateStr func(string) string
	var replaceStartPos int
	var name string

	// %ENVNAME%
	percentCount := strings.Count(cmdline, "%")
	if percentCount%2 == 1 {
		replaceStartPos = strings.LastIndex(cmdline, "%")
		if replaceStartPos < 0 {
			return nil, -1, nil
		}
		name = cmdline[replaceStartPos+1:]
		makeCandidateStr = func(name string) string {
			return "%" + name + "%"
		}
	} else {
		// $ENVNAME
		m := rxDollar.FindStringIndex(cmdline)
		if len(m) > 0 {
			replaceStartPos = m[0]
			name = cmdline[m[0]+1:]
			makeCandidateStr = func(name string) string {
				return "$" + name
			}
		} else {
			m = rxDollar2.FindStringIndex(cmdline)
			if m == nil || len(m) <= 0 {
				return nil, -1, nil
			}
			replaceStartPos = m[0]
			name = cmdline[m[0]+2:]
			makeCandidateStr = func(name string) string {
				return "${" + name + "}"
			}
		}
	}

	name = strings.ToUpper(name)
	var matches []Element

	for _, vars := range PercentVariables {
		vars.EachKey(func(envName string) {
			if strings.HasPrefix(strings.ToUpper(envName), name) {
				envValue := makeCandidateStr(envName)
				element := Element1(envValue)
				matches = append(matches, element)
			}
		})
	}
	if len(matches) <= 0 { // nothing matches.
		return nil, -1, nil
	}
	return matches, replaceStartPos, nil
}
