package completion

import (
	"context"
	"os"
	"strings"
)

func completionSet(ctx context.Context, params []string) ([]Element, error) {
	result := []Element{}
	base := strings.ToUpper(params[len(params)-1])
	for _, env1 := range os.Environ() {
		if strings.HasPrefix(strings.ToUpper(env1), base) {
			result = append(result, Element1(env1))
		}
	}
	return result, nil
}

func completionCd(ctx context.Context, params []string) ([]Element, error) {
	return listUpDirs(ctx, params[len(params)-1])
}
