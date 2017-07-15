package mains

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/zetamatta/nyagos/dos"
	. "github.com/zetamatta/nyagos/ifdbg"
	"github.com/zetamatta/nyagos/lua"
	"github.com/zetamatta/nyagos/shell"
)

func versionOrStamp() string {
	if Version != "" {
		return Version
	} else {
		return "v" + Stamp
	}
}

func loadBundleScript1(it *shell.Cmd, L lua.Lua, path string) error {
	if DBG {
		println("load cached ", path)
	}
	bin, err := Asset(path)
	if err != nil {
		return err
	}
	err = L.LoadBufferX(path, bin, "t")
	if err != nil {
		return err
	}
	err = NyagosCallLua(L, it, 0, 0)
	if err != nil {
		return err
	}
	return nil
}

type InterpreterT interface {
	Interpret(string) (int, error)
}

func loadScripts(it *shell.Cmd, L lua.Lua) error {
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
		fmt.Println("(c) 2014-2017 NYAOS.ORG <http://www.nyaos.org>")
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
				if !strings.HasSuffix(strings.ToLower(name1), ".lua") {
					continue
				}
				relpath := "nyagos.d/" + name1
				asset1, assetErr := AssetInfo(relpath)
				if assetErr == nil && asset1.Size() == finfo1.Size() && !asset1.ModTime().Truncate(time.Second).Before(finfo1.ModTime().Truncate(time.Second)) {
					if err := loadBundleScript1(it, L, relpath); err != nil {
						fmt.Fprintf(os.Stderr, "cached %s: %s\n", relpath, err)
					}
				} else {
					path1 := filepath.Join(nyagos_d, name1)
					if DBG {
						println("load real ", path1)
					}
					if err := L.Source(path1); err != nil {
						fmt.Fprintf(os.Stderr, "%s: %s\n", name1, err.Error())
					}
				}
			}
		}
	} else if assertdir, err := AssetDir("nyagos.d"); err == nil {
		// nyagos.d/ not found.
		for _, name1 := range assertdir {
			if !strings.HasSuffix(strings.ToLower(name1), ".lua") {
				continue
			}
			relpath := "nyagos.d/" + name1
			if err1 := loadBundleScript1(it, L, relpath); err1 != nil {
				fmt.Fprintf(os.Stderr, "bundled %s: %s\n", relpath, err1.Error())
			}
		}
	}
	if _, err := runLua(it, L, filepath.Join(exeFolder, ".nyagos")); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	barNyagos(it, exeFolder, L)
	if err := dotNyagos(it, L); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	barNyagos(it, dos.GetHome(), L)
	return nil
}

func runLua(it *shell.Cmd, L lua.Lua, fname string) ([]byte, error) {
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
	if err := NyagosCallLua(L, it, 0, 0); err != nil {
		return nil, err
	}
	// println("Run: " + fname)
	return chank, nil
}

func dotNyagos(it *shell.Cmd, L lua.Lua) error {
	dot_nyagos := filepath.Join(dos.GetHome(), ".nyagos")
	dotStat, err := os.Stat(dot_nyagos)
	if err != nil {
		return nil
	}
	cachePath := filepath.Join(AppDataDir(), runtime.GOARCH+".nyagos.luac")
	cacheStat, err := os.Stat(cachePath)
	if err == nil && !dotStat.ModTime().After(cacheStat.ModTime()) {
		_, err = runLua(it, L, cachePath)
		return err
	}
	chank, err := runLua(it, L, dot_nyagos)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(cachePath, chank, os.FileMode(0644))
}

func barNyagos(it InterpreterT, folder string, L lua.Lua) {
	bar_nyagos := filepath.Join(folder, "_nyagos")
	fd, fd_err := os.Open(bar_nyagos)
	if fd_err != nil {
		return
	}
	defer fd.Close()
	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		text := scanner.Text()
		text = doLuaFilter(L, text)
		_, err := it.Interpret(text)
		if err != nil {
			fmt.Fprint(os.Stderr, err.Error())
		}
	}
}
