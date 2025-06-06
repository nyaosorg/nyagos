package frame

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
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

type scriptEngine interface {
	DoFile(string) error
	DoString(string) error
}

func loadScriptDir(dir string, L scriptEngine) error {
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

		if strings.HasSuffix(lowerName, ".lua") {
			err := L.DoFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", path, err.Error())
			}
		}
	}
	return nil
}

func loadEmbedScripts(L scriptEngine, fsys fs.FS, dir string, warn func(error) error) error {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		return err
	}
	for _, entry1 := range entries {
		scriptName := entry1.Name()
		if len(scriptName) > 0 && scriptName[0] == '.' {
			continue
		}
		scriptPath := path.Join(dir, scriptName)
		if entry1.IsDir() {
			if err := loadEmbedScripts(L, fsys, scriptPath, warn); err != nil {
				return err
			}
			continue
		}
		if !strings.HasSuffix(scriptName, ".lua") {
			continue
		}
		scriptCode, err := fs.ReadFile(fsys, scriptPath)
		if err != nil {
			return err
		}
		if err := L.DoString(string(scriptCode)); err != nil {
			if err = warn(fmt.Errorf("(embed) %s: %w", scriptName, err)); err != nil {
				return err
			}
		}
	}
	return nil
}
func LoadEmbedScripts(L scriptEngine, fsys fs.FS, warn func(error) error) error {
	return loadEmbedScripts(L, fsys, ".", warn)
}

// LoadScripts loads ".nyagos"
func LoadScripts(L scriptEngine, warn func(error) error) error {
	exeName, err := os.Executable()
	if err != nil {
		exeName = os.Args[0]
	}
	exeFolder := filepath.Dir(exeName)
	loadScriptDir(filepath.Join(exeFolder, "nyagos.d"), L)

	if appDir, err := os.UserConfigDir(); err == nil {
		dir := filepath.Join(appDir, "NYAOS_ORG/nyagos.d")
		err := loadScriptDir(dir, L)
		if err != nil {
			if _, ok := err.(_DirNotFound); ok {
				os.MkdirAll(dir, 0755)
			} else {
				if err := warn(err); err != nil {
					return err
				}
			}
		}
	}

	fname := filepath.Join(exeFolder, ".nyagos")
	if _, err := os.Stat(fname); err == nil {
		if err := L.DoFile(fname); err != nil {
			if err = warn(err); err != nil {
				return err
			}
		}
	}
	if err := dotNyagos(L); err != nil {
		if err = warn(err); err != nil {
			return err
		}
	}
	return nil
}

func dotNyagos(L scriptEngine) error {
	dotNyagos := filepath.Join(nodos.GetHome(), ".nyagos")
	if _, err := os.Stat(dotNyagos); err != nil {
		return nil
	}
	return L.DoFile(dotNyagos)
}
