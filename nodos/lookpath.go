package nodos

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/zetamatta/go-findfile"
)

func lookPath(dir1, patternBase string) (foundpath string) {
	pathExtList := filepath.SplitList(os.Getenv("PATHEXT"))
	ext := filepath.Ext(patternBase)
	var pattern string
	if ext == "" {
		pattern = patternBase + ".*"
	} else {
		pattern = patternBase
	}
	basename := filepath.Base(patternBase)
	names := map[string]int{strings.ToUpper(basename): 0}
	for i, ext1 := range pathExtList {
		names[strings.ToUpper(basename+ext1)] = i + 1
	}
	foundIndex := 999
	findfile.Walk(pattern, func(f *findfile.FileInfo) bool {
		if f.IsDir() {
			return true
		}
		if filepath.Ext(f.Name()) == "" {
			return true
		}
		if i, ok := names[strings.ToUpper(f.Name())]; ok && i < foundIndex {
			foundIndex = i
			foundpath = filepath.Join(dir1, f.Name())
			if f.IsReparsePoint() {
				var err error
				linkTo, err := os.Readlink(foundpath)
				if err == nil && linkTo != "" {
					foundpath = linkTo
					if !filepath.IsAbs(foundpath) {
						foundpath = filepath.Join(dir1, foundpath)
					}
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
