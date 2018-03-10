package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zetamatta/nyagos/alias"
	"github.com/zetamatta/nyagos/dos"
)

const (
	WHICH_NOT_FOUND = 1
)

func envToList(first1 string, envs ...string) []string {
	result := make([]string, 1, 20)
	result[0] = first1
	for _, env := range envs {
		list1 := filepath.SplitList(os.Getenv(env))
		result = append(result, list1...)
	}
	return result
}

func cmdWhich(ctx context.Context, cmd Param) (int, error) {
	all := false
	var pathList []string
	var extList []string
	for _, name := range cmd.Args()[1:] {
		if name == "-a" {
			all = true
			pathList = envToList(".", "PATH", "NYAGOSPATH")
			extList = envToList("", "PATHEXT")
			continue
		}
		if a, ok := alias.Table[strings.ToLower(name)]; ok {
			fmt.Fprintf(cmd.Out(), "%s: aliased to %s\n", name, a.String())
			if !all {
				continue
			}
		}
		if _, ok := buildInCommand[name]; ok {
			fmt.Fprintf(cmd.Out(), "%s: built-in command\n", name)
			if !all {
				continue
			}
		}
		if all {
			for _, dir1 := range pathList {
				for _, ext1 := range extList {
					fullpath1 := filepath.Join(dir1, name)
					fullpath1 = fullpath1 + ext1
					if _, err1 := os.Stat(fullpath1); err1 == nil {
						fmt.Fprintln(cmd.Out(), fullpath1)
					}
				}
			}

		} else {
			path := dos.LookPath(name, "NYAGOSPATH")
			if path == "" {
				return WHICH_NOT_FOUND, fmt.Errorf("which %s: not found", name)
			}
			fmt.Fprintln(cmd.Out(), filepath.Clean(path))
		}
	}
	return 0, nil
}
