package completion

import (
	"os"
	"path"
	"regexp"
	"strings"

	"../dos"
)

var rxEnvPattern = regexp.MustCompile("%[^%]+%")
var rxTilde = regexp.MustCompile("^~[/\\\\]")

func listUpFiles(str string) ([]string, error) {
	str = rxEnvPattern.ReplaceAllStringFunc(str, func(p string) string {
		if len(p) == 2 {
			return "%"
		} else if val := os.Getenv(p[1 : len(p)-1]); val != "" {
			return val
		} else {
			return p
		}
	})

	str = rxTilde.ReplaceAllStringFunc(str, func(p string) string {
		if home := dos.GetHome(); home != "" {
			return home + "\\"
		} else {
			return p
		}
	})
	str = strings.Replace(strings.Replace(str, "\\", "/", -1), "\"", "", -1)
	var directory string
	str = strings.Replace(str, "\\", "/", -1)
	if strings.HasSuffix(str, "/") {
		directory = path.Join(str, ".")
	} else {
		directory = path.Dir(str)
	}
	cutprefix := 0
	if strings.HasPrefix(directory, "/") {
		wd, _ := os.Getwd()
		directory = wd[0:2] + directory
		cutprefix = 2
	}

	fd, fdErr := os.Open(directory)
	if fdErr != nil {
		return nil, fdErr
	}
	defer fd.Close()
	files, filesErr := fd.Readdir(-1)
	if filesErr != nil {
		return nil, filesErr
	}
	commons := make([]string, 0)
	STR := strings.ToUpper(str)
	for _, node1 := range files {
		name := path.Join(directory, node1.Name())
		if attr, attrErr := dos.GetFileAttributes(name); attrErr == nil && (attr&dos.FILE_ATTRIBUTE_HIDDEN) != 0 {
			continue
		}
		if node1.IsDir() {
			name += "/"
		}
		if cutprefix > 0 {
			name = name[2:]
		}
		nameUpr := strings.ToUpper(name)
		if strings.HasPrefix(nameUpr, STR) {
			commons = append(commons, name)
		}
	}
	return commons, nil
}
