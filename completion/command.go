package completion

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"../cpath"
)

func isExecutable(path string) bool {
	return cpath.IsExecutableSuffix(filepath.Ext(path))
}

func listUpAllExecutableOnEnv(envName string) []string {
	list := make([]string, 0, 100)
	pathEnv := os.Getenv(envName)
	dirList := filepath.SplitList(pathEnv)
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
			if file1.IsDir() {
				continue
			}
			name := file1.Name()
			if isExecutable(name) {
				list = append(list, path.Base(name))
			}
		}
	}
	return list
}

func listUpCurrentAllExecutable(str string) ([]string, error) {
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
	return list, nil
}

func removeDup(list []string) []string {
	found := map[string]bool{}
	result := make([]string, 0, len(list))

	for _, value := range list {
		if _, ok := found[value]; !ok {
			result = append(result, value)
			found[value] = true
		}
	}
	return result
}

func listUpCommands(str string) ([]string, error) {
	list, listErr := listUpCurrentAllExecutable(str)
	if listErr != nil {
		return nil, listErr
	}
	strUpr := strings.ToUpper(str)
	for _, f := range command_listupper {
		for _, name := range f() {
			name1Upr := strings.ToUpper(name)
			if strings.HasPrefix(name1Upr, strUpr) {
				list = append(list, name)
			}
		}
	}
	return removeDup(list), nil
}
