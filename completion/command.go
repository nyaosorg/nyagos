package completion

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func listUpAllExecutableOnEnv(ctx context.Context, envName string) ([]Element, error) {
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
		files, err := ioutil.ReadDir(dir1)
		if err != nil {
			continue
		}
		for _, file1 := range files {
			if file1.IsDir() {
				continue
			}
			name := file1.Name()
			if isExecutable(name) {
				name_ := path.Base(name)
				element := Element1(name_)
				list = append(list, element)
			}
		}
	}
	return list, nil
}

func listUpCurrentAllExecutable(ctx context.Context, str string) ([]Element, error) {
	listTmp, listErr := listUpFiles(ctx, str)
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

func listUpCommands(ctx context.Context, str string) ([]Element, error) {
	list, listErr := listUpCurrentAllExecutable(ctx, str)
	if listErr != nil {
		return nil, listErr
	}
	strUpr := strings.ToUpper(str)
	for _, f := range commandListUpper {
		files, err := f(ctx)
		if err != nil {
			return nil, err
		}
		for _, element := range files {
			name1Upr := strings.ToUpper(element.String())
			if strings.HasPrefix(name1Upr, strUpr) {
				list = append(list, element)
			}
		}
	}
	return removeDup(list), nil
}
