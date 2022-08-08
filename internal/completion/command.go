package completion

import (
	"context"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func listUpAllFilesOnEnv(ctx context.Context, envName string, filter func(fs.DirEntry) bool) ([]Element, error) {
	list := make([]Element, 0, 100)
	pathEnv := os.Getenv(envName)
	dirList := filepath.SplitList(pathEnv)
	for _, dir1 := range dirList {
		if ctx != nil {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}
		files, err := os.ReadDir(dir1)
		if err != nil {
			continue
		}
		for _, file1 := range files {
			if filter(file1) {
				name := file1.Name()
				_name := path.Base(name)
				element := Element1(_name)
				list = append(list, element)
			}
		}
	}
	return list, nil
}

func listUpAllExecutableOnEnv(ctx context.Context, envName string) ([]Element, error) {
	return listUpAllFilesOnEnv(ctx, envName, func(file1 fs.DirEntry) bool {
		return !file1.IsDir() && isExecutable(file1.Name())
	})
}

func listUpCurrentAllExecutable(ctx context.Context, str string) ([]Element, error) {
	listTmp, listErr := ListUpFiles(ctx, DoNotUncCompletion, str)
	if listErr != nil {
		return nil, listErr
	}
	list := make([]Element, 0, len(listTmp))
	for _, p := range listTmp {
		if endWithRoot(p.String()) || isExecutable(p.String()) {
			list = append(list, p)
		}
	}
	return list, nil
}

func removeDup(list []Element) []Element {
	found := map[string]struct{}{}
	result := make([]Element, 0, len(list))

	for _, value := range list {
		if _, ok := found[value.String()]; !ok {
			result = append(result, value)
			found[value.String()] = struct{}{}
		}
	}
	return result
}

func filterElementWithPrefix(dest, source []Element, prefix string) []Element {
	upperPrefix := strings.ToUpper(prefix)
	for _, element := range source {
		name1Upr := strings.ToUpper(element.String())
		if strings.HasPrefix(name1Upr, upperPrefix) {
			dest = append(dest, element)
		}
	}
	return dest
}

func listUpCommands(ctx context.Context, str string) ([]Element, error) {
	list, listErr := listUpCurrentAllExecutable(ctx, str)
	if listErr != nil {
		return nil, listErr
	}
	for _, f := range commandListUpper {
		files, err := f(ctx)
		if err != nil {
			return nil, err
		}
		list = filterElementWithPrefix(list, files, str)
	}
	return removeDup(list), nil
}
