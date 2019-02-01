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

func completionEnv(ctx context.Context, param []string) ([]Element, error) {
	eq := -1
	for i := 1; i < len(param); i++ {
		if strings.Contains(param[i], "=") {
			eq = i
		}
	}
	current := len(param) - 1

	if current == eq || current == 1 {
		return completionSet(ctx, param)
	} else if current == eq+1 {
		return listUpCommands(ctx, param[current])
	} else {
		return nil, nil
	}
}
