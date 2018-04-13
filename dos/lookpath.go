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
		for _, name1 := range names {
			if strings.EqualFold(f.Name(), name1) {
				foundpath = filepath.Join(dir1, f.Name())
				if !f.IsReparsePoint() {
					return false
				}
				var err error
				foundpath_, err := os.Readlink(foundpath)
				if err == nil {
					if foundpath_ != "" {
						foundpath = foundpath_
						if filepath.IsAbs(foundpath) {
							return false
						}
						foundpath = filepath.Join(dir1, foundpath)
						return false
					}
				} else if defined.DBG {
					print(err.Error(), "\n")
				}
			}
		}
		return true
	})
	return
}

func LookPath(name string, envnames ...string) string {
	if strings.ContainsAny(name, "\\/:") {
		return lookPath(filepath.Dir(name), name)
	}
	var envlist strings.Builder
	envlist.WriteRune('.')
	envlist.WriteRune(os.PathListSeparator)
	envlist.WriteString(os.Getenv("PATH"))
	for _, name1 := range envnames {
		envlist.WriteRune(os.PathListSeparator)
		envlist.WriteString(os.Getenv(name1))
	}
	// println(envlist.String())
	pathDirList := filepath.SplitList(envlist.String())

	for _, dir1 := range pathDirList {
		// println("lookPath:" + dir1)
		if path := lookPath(dir1, filepath.Join(dir1, name)); path != "" {
			// println("Found:" + path)
			return path
		}
	}
	return ""
}
