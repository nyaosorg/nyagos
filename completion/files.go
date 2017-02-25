package completion

import (
	"os"
	"regexp"
	"strings"

	"github.com/zetamatta/go-findfile"

	"../cpath"
)

const (
	STD_SLASH = string(os.PathSeparator)
	OPT_SLASH = "/"
)

var rxEnvPattern = regexp.MustCompile("%[^%]+%")
var rxTilde = regexp.MustCompile("^~[/\\\\]")

func replaceEnv(str string) string {
	str = rxEnvPattern.ReplaceAllStringFunc(str, func(p string) string {
		if len(p) == 2 {
			return "%"
		}
		name := p[1 : len(p)-1]
		for _, env := range PercentVariables {
			if value := env.Lookup(name); value != "" {
				return value
			}
		}
		return p
	})

	str = rxTilde.ReplaceAllStringFunc(str, func(p string) string {
		if home := cpath.GetHome(); home != "" {
			return home + "\\"
		} else {
			return p
		}
	})

	return str
}

func listUpFiles(str string) ([]Element, error) {
	orgSlash := STD_SLASH[0]
	if pos := strings.IndexAny(str, STD_SLASH+OPT_SLASH); pos >= 0 {
		orgSlash = str[pos]
	}
	str = strings.Replace(strings.Replace(str, OPT_SLASH, STD_SLASH, -1), "\"", "", -1)
	directory := cpath.DirName(str)
	wildcard := cpath.Join(replaceEnv(directory), "*")

	// Drive letter
	cutprefix := 0
	if strings.HasPrefix(directory, STD_SLASH) {
		wd, _ := os.Getwd()
		directory = wd[0:2] + directory
		cutprefix = 2
	}
	commons := make([]Element, 0)
	STR := strings.ToUpper(str)
	fdErr := findfile.Walk(wildcard, func(fd *findfile.FileInfo) bool {
		if fd.Name() == "." || fd.Name() == ".." || fd.IsHidden() {
			return true
		}
		listname := fd.Name()
		name := cpath.Join(directory, fd.Name())
		if fd.IsDir() {
			name += STD_SLASH
			listname += OPT_SLASH
		}
		if cutprefix > 0 {
			name = name[2:]
		}
		nameUpr := strings.ToUpper(name)
		if strings.HasPrefix(nameUpr, STR) {
			if orgSlash != STD_SLASH[0] {
				name = strings.Replace(name, STD_SLASH, OPT_SLASH, -1)
			}
			element := Element{InsertStr: name, ListupStr: listname}
			commons = append(commons, element)
		}
		return true
	})
	return commons, fdErr
}
