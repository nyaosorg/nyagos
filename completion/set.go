package completion

import (
	"context"
	"os"
	"strings"

	"github.com/mitchellh/go-ps"
)

func completionSet(ctx context.Context, ua UncAccess, params []string) ([]Element, error) {
	result := []Element{}
	base := strings.ToUpper(params[len(params)-1])
	for _, env1 := range os.Environ() {
		if strings.HasPrefix(strings.ToUpper(env1), base) {
			result = append(result, Element1(env1))
		}
	}
	return result, nil
}

func completionCd(ctx context.Context, ua UncAccess, params []string) ([]Element, error) {
	return listUpDirs(ctx, ua, params[len(params)-1])
}

func completionEnv(ctx context.Context, ua UncAccess, param []string) ([]Element, error) {
	eq := -1
	for i := 1; i < len(param); i++ {
		if strings.Contains(param[i], "=") {
			eq = i
		}
	}
	current := len(param) - 1

	if current == eq || current == 1 {
		return completionSet(ctx, ua, param)
	} else if current == eq+1 {
		return listUpCommands(ctx, param[current])
	} else {
		return nil, nil
	}
}

func completionWhich(ctx context.Context, ua UncAccess, param []string) ([]Element, error) {
	if len(param) == 2 {
		return listUpCommands(ctx, param[len(param)-1])
	}
	return nil, nil
}

func completionProcessName(ctx context.Context, ua UncAccess, param []string) ([]Element, error) {
	processes, err := ps.Processes()
	if err != nil {
		return nil, err
	}
	uniq := map[string]struct{}{}
	base := strings.ToUpper(param[len(param)-1])
	for _, ps1 := range processes {
		name := ps1.Executable()
		if strings.HasPrefix(strings.ToUpper(name), base) {
			uniq[name] = struct{}{}
		}
	}
	result := make([]Element, 0, len(uniq))
	for name := range uniq {
		result = append(result, Element1(name))
	}
	return result, nil
}

func completionTaskKill(ctx context.Context, ua UncAccess, param []string) ([]Element, error) {
	if len(param) >= 3 && strings.EqualFold(param[len(param)-2], "/IM") {
		return completionProcessName(ctx, ua, param)
	}
	return nil, nil
}
