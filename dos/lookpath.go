package dos

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zetamatta/go-findfile"

	"github.com/zetamatta/nyagos/defined"
)

func lookPath(dir1, patternBase string) (foundpath string) {
	pattern := patternBase + ".*"
	pathExtList := filepath.SplitList(os.Getenv("PATHEXT"))
	names := make([]string, len(pathExtList)+1)
	basename := filepath.Base(patternBase)
	names[0] = basename
	for i, ext1 := range pathExtList {
		names[i+1] = basename + ext1
	}
	findfile.Walk(pattern, func(f *findfile.FileInfo) bool {
		if f.IsDir() {
			return true
		}
		if filepath.Ext(f.Name()) == "" {
			return true
		}
		for _, name1 := range names {
			if strings.EqualFold(f.Name(), name1) {
				foundpath = filepath.Join(dir1, f.Name())
				if !f.IsReparsePoint() {
					return false
				}
				var err error
				linkTo, err := os.Readlink(foundpath)
				if err == nil && linkTo != "" {
					foundpath = linkTo
					if filepath.IsAbs(foundpath) {
						return false
					}
					foundpath = filepath.Join(dir1, foundpath)
					return false
				} else if defined.DBG {
					print(err.Error(), "\n")
				}
			}
		}
		return true
	})
	return
}

// LookCurdirT is the type for constant meaning the current directory should be looked.
type LookCurdirT int

const (
	// LookCurdirFirst means that the current directory should be looked at first.
	LookCurdirFirst LookCurdirT = iota
	// LookCurdirLast  means that the current directory should be looked at last.
	LookCurdirLast
	// LookCurdirNever menas that the current directory should be never looked.
	LookCurdirNever
)

// LookPath search `name` from %PATH% and the directories listed by
// the environment variables `envnames`.
func LookPath(where LookCurdirT, name string, envnames ...string) string {
	if strings.ContainsAny(name, "\\/:") {
		return lookPath(filepath.Dir(name), name)
	}
	var envlist strings.Builder
	if where == LookCurdirFirst {
		envlist.WriteRune('.')
		envlist.WriteRune(os.PathListSeparator)
	}
	envlist.WriteString(os.Getenv("PATH"))
	if where == LookCurdirLast {
		envlist.WriteRune(os.PathListSeparator)
		envlist.WriteRune('.')
	}
	for _, name1 := range envnames {
		envlist.WriteRune(os.PathListSeparator)
		envlist.WriteString(os.Getenv(name1))
	}
	// println(envlist.String())
	pathDirList := filepath.SplitList(envlist.String())

	for _, dir1 := range pathDirList {
		// println("lookPath:" + dir1)
		_dir1 := strings.TrimSpace(dir1)
		if _dir1 == "" {
			continue
		}
		if path := lookPath(dir1, filepath.Join(_dir1, name)); path != "" {
			// println("Found:" + path)
			return path
		}
	}
	return ""
}
