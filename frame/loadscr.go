package frame

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zetamatta/nyagos/nodos"
)

// Version is to show title display.
var Version string

type DirNotFound struct {
	err error
}

func (e DirNotFound) Error() string {
	return e.err.Error()
}

func (e DirNotFound) Unwrap() error {
	return e.err
}

func loadScriptDir(dir string,
	shellEngine func(string) error,
	langEngine func(string) ([]byte, error)) error {

	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return DirNotFound{err: err}
		}
		return err
	}

	for _, f := range files {
		name := f.Name()
		path := filepath.Join(dir, name)
		lowerName := strings.ToLower(name)

		var err error
		if strings.HasSuffix(lowerName, ".lua") {
			_, err = langEngine(path)
		} else if strings.HasSuffix(lowerName, ".ny") {
			err = shellEngine(path)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", path, err.Error())
		}
	}
	return nil
}

func LoadScripts(
	shellEngine func(string) error,
	langEngine func(string) ([]byte, error)) error {

	exeName, err := os.Executable()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	exeFolder := filepath.Dir(exeName)
	loadScriptDir(filepath.Join(exeFolder, "nyagos.d"),
		shellEngine, langEngine)

	if appDir, err := os.UserConfigDir(); err == nil {
		dir := filepath.Join(appDir, "NYAOS_ORG/nyagos.d")
		err := loadScriptDir(dir, shellEngine, langEngine)
		if err != nil {
			if _, ok := err.(DirNotFound); ok {
				os.MkdirAll(dir, 0755)
			} else {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}

	fname := filepath.Join(exeFolder, ".nyagos")
	if _, err := os.Stat(fname); err == nil {
		if _, err := langEngine(fname); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	barNyagos(shellEngine, exeFolder)
	if err := dotNyagos(langEngine); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	barNyagos(shellEngine, nodos.GetHome())
	return nil
}

func dotNyagos(langEngine func(string) ([]byte, error)) error {
	dot_nyagos := filepath.Join(nodos.GetHome(), ".nyagos")
	dotStat, err := os.Stat(dot_nyagos)
	if err != nil {
		return nil
	}
	cachePath := filepath.Join(AppDataDir(), runtime.GOARCH+".nyagos.luac")
	cacheStat, err := os.Stat(cachePath)
	if err == nil {
		if cacheStat.Size() != 0 && !dotStat.ModTime().After(cacheStat.ModTime()) {
			_, err = langEngine(cachePath)
			if err == nil {
				return nil
			}
		}
		os.Remove(cachePath)
	}
	chank, err := langEngine(dot_nyagos)
	if err != nil || chank == nil {
		return err
	}
	return os.WriteFile(cachePath, chank, os.FileMode(0644))
}

func barNyagos(shellEngine func(string) error, folder string) {
	bar_nyagos := filepath.Join(folder, "_nyagos")
	fd, err := os.Open(bar_nyagos)
	if err != nil {
		return
	}
	err = shellEngine(bar_nyagos)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
	}
	fd.Close()
}
