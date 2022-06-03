package frame

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

// Version is to show title display.
var Version string

type _DirNotFound struct {
	err error
}

func (e _DirNotFound) Error() string {
	return e.err.Error()
}

func (e _DirNotFound) Unwrap() error {
	return e.err
}

func loadScriptDir(dir string,
	shellEngine func(string) error,
	langEngine func(string) ([]byte, error)) error {

	files, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return _DirNotFound{err: err}
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

// LoadScripts loads ".nyagos" an "_nyagos"
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
			if _, ok := err.(_DirNotFound); ok {
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
	dotNyagos := filepath.Join(nodos.GetHome(), ".nyagos")
	dotStat, err := os.Stat(dotNyagos)
	if err != nil {
		return nil
	}
	cachePath := filepath.Join(appDataDir(), runtime.GOARCH+".nyagos.luac")
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
	chank, err := langEngine(dotNyagos)
	if err != nil || chank == nil {
		return err
	}
	return os.WriteFile(cachePath, chank, os.FileMode(0644))
}

func barNyagos(shellEngine func(string) error, folder string) {
	barNyagos := filepath.Join(folder, "_nyagos")
	fd, err := os.Open(barNyagos)
	if err != nil {
		return
	}
	err = shellEngine(barNyagos)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
	}
	fd.Close()
}
