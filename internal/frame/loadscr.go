package frame

import (
	"fmt"
	"os"
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

// LoadScripts loads ".nyagos"
func LoadScripts(L scriptEngine) error {
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
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}

	fname := filepath.Join(exeFolder, ".nyagos")
	if _, err := os.Stat(fname); err == nil {
		if err := L.DoFile(fname); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
	if err := dotNyagos(L); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
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
