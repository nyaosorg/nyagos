package completion

import (
	"os"
	"path"
	"strings"
)

func listUpCommands(str string) ([]string, error) {
	listTmp, listErr := listUpFiles(str)
	if listErr != nil {
		return nil, listErr
	}
	list := make([]string, 0)
	for _, fname := range listTmp {
		if strings.HasSuffix(fname, "/") || strings.HasSuffix(fname, "\\") || isExecutable(fname) {
			list = append(list, fname)
		}
	}
	pathEnv := os.Getenv("PATH")
	dirList := strings.Split(pathEnv, ";")
	strUpr := strings.ToUpper(str)
	for _, dir1 := range dirList {
		dirHandle, dirErr := os.Open(dir1)
		if dirErr != nil {
			continue
		}
		defer dirHandle.Close()
		files, filesErr := dirHandle.Readdir(0)
		if filesErr != nil {
			continue
		}
		for _, file1 := range files {
			name1Upr := strings.ToUpper(file1.Name())
			if !strings.HasPrefix(name1Upr, strUpr) {
				continue
			}
			if file1.IsDir() {
				continue
			}
			name := file1.Name()
			if isExecutable(name) {
				list = append(list, path.Base(name))
			}
		}
	}
	// remove dupcalites
	uniq := make([]string, 0)
	lastone := ""
	for _, cur := range list {
		if cur != lastone {
			uniq = append(uniq, cur)
		}
		lastone = cur
	}
	return uniq, nil
}
