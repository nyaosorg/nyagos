package mains

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/zetamatta/nyagos/dos"
	"github.com/zetamatta/nyagos/lua"
	"github.com/zetamatta/nyagos/shell"
)

func versionOrStamp() string {
	if Version != "" {
		return Version
	} else if Stamp != "" {
		return "v" + Stamp
	} else {
		return "snapshot"
	}
}

func loadScripts(shellEngine func(string) error,
	langEngine func(string) ([]byte, error)) error {

	exeName, exeNameErr := os.Executable()
	if exeNameErr != nil {
		fmt.Fprintln(os.Stderr, exeNameErr)
	}
	exeFolder := filepath.Dir(exeName)

	if !silentmode {
		fmt.Printf("Nihongo Yet Another GOing Shell %s-%s by %s & %s\n",
			versionOrStamp(),
			runtime.GOARCH,
			runtime.Version(),
			"Lua 5.3")
		fmt.Println("(c) 2014-2018 NYAOS.ORG <http://www.nyaos.org>")
	}

	nyagos_d := filepath.Join(exeFolder, "nyagos.d")
	nyagos_d_fd, nyagos_d_err := os.Open(nyagos_d)
	if nyagos_d_err == nil {
		defer nyagos_d_fd.Close()
		finfos, err := nyagos_d_fd.Readdir(-1)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} else {
			for _, finfo1 := range finfos {
				name1 := finfo1.Name()
				path1 := filepath.Join(nyagos_d, name1)
				name1_ := strings.ToLower(name1)

				var err error
				if strings.HasSuffix(name1_, ".lua") {
					_, err = langEngine(path1)
				} else if strings.HasSuffix(name1_, ".ny") {
					err = shellEngine(path1)
				}
				if err != nil {
					fmt.Fprintf(os.Stderr, "%s: %s\n", name1, err.Error())
				}
			}
		}
	}
	if _, err := langEngine(filepath.Join(exeFolder, ".nyagos")); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	barNyagos(shellEngine, exeFolder)
	if err := dotNyagos(langEngine); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	barNyagos(shellEngine, dos.GetHome())
	return nil
}

func runLua(it *shell.Shell, L lua.Lua, fname string) ([]byte, error) {
	_, err := os.Stat(fname)
	if err != nil {
		if os.IsNotExist(err) {
			// println("pass " + fname + " (not exists)")
			return []byte{}, nil
		} else {
			return nil, err
		}
	}
	if _, err := L.LoadFile(fname, "bt"); err != nil {
		return nil, err
	}
	chank := L.Dump()
	if err := callLua(it, 0, 0); err != nil {
		return nil, err
	}
	// println("Run: " + fname)
	return chank, nil
}

func dotNyagos(langEngine func(string) ([]byte, error)) error {
	dot_nyagos := filepath.Join(dos.GetHome(), ".nyagos")
	dotStat, err := os.Stat(dot_nyagos)
	if err != nil {
		return nil
	}
	cachePath := filepath.Join(AppDataDir(), runtime.GOARCH+".nyagos.luac")
	cacheStat, err := os.Stat(cachePath)
	if err == nil && !dotStat.ModTime().After(cacheStat.ModTime()) {
		_, err = langEngine(cachePath)
		return err
	}
	chank, err := langEngine(dot_nyagos)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cachePath, chank, os.FileMode(0644))
}

func barNyagos(shellEngine func(string) error, folder string) {
	bar_nyagos := filepath.Join(folder, "_nyagos")
	fd, err := os.Open(bar_nyagos)
	if err != nil {
		return
	}
	err = shellEngine(bar_nyagos)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
	}
	fd.Close()
}
