package frame

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

// Version is to show title display.
var Version string

type scriptEngine interface {
	DoFile(string) error
	DoString(string) error
}

func LoadScriptsFs(L scriptEngine, fsys fs.FS, warn func(error) error) error {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return err
	}
	for _, entry1 := range entries {
		if entry1.IsDir() {
			continue
		}
		scriptName := entry1.Name()
		if !strings.HasSuffix(scriptName, ".lua") {
			continue
		}
		scriptCode, err := fs.ReadFile(fsys, scriptName)
		if err != nil {
			return err
		}
		if err := L.DoString(string(scriptCode)); err != nil {
			if err = warn(fmt.Errorf("%s: %w", scriptName, err)); err != nil {
				return err
			}
		}
	}
	return nil
}

func LoadScripts(L scriptEngine, warn func(error) error) error {
	exeName, err := os.Executable()
	if err != nil {
		exeName = os.Args[0]
	}
	exeFolder := filepath.Dir(exeName)
	nyagosD := filepath.Join(exeFolder, "nyagos.d")
	if err := LoadScriptsFs(L, os.DirFS(nyagosD), warn); err != nil {
		return err
	}
	if appDir, err := os.UserConfigDir(); err == nil {
		dir := filepath.Join(appDir, "NYAOS_ORG/nyagos.d")
		if _, err := os.Stat(dir); err != nil {
			os.MkdirAll(dir, 0755)
		} else if err := LoadScriptsFs(L, os.DirFS(dir), warn); err != nil {
			return err
		}
	}
	dirs := []string{exeFolder, nodos.GetHome()}
	for _, dir1 := range dirs {
		fname := filepath.Join(dir1, ".nyagos")
		if _, err := os.Stat(fname); err == nil {
			if err := L.DoFile(fname); err != nil {
				if err = warn(err); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
