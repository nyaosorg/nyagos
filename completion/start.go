package completion

import (
	"context"
	"io/fs"
	"strings"
)

func completionStart(ctx context.Context, ua UncCompletion, param []string) ([]Element, error) {
	if len(param) <= 0 {
		return []Element{}, nil
	}
	baseName := param[len(param)-1]

	elements, err := ListUpFiles(ctx, ua, baseName)
	if err != nil {
		return nil, err
	}
	elements = filterElementWithPrefix(nil, elements, baseName)
	if !strings.ContainsAny(baseName, `\/`) {
		_elements, err := listUpAllFilesOnEnv(ctx,
			"PATH",
			func(fs.DirEntry) bool { return true })
		if err == nil {
			elements = filterElementWithPrefix(elements, _elements, baseName)
		}
	}
	return removeDup(elements), nil
}
