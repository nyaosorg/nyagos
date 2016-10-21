package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"../dos"
	"../lua"
)

func versionOrStamp() string {
	if version != "" {
		return version
	} else {
		return "v" + stamp
	}
}

func loadBundleScript1(L lua.Lua, path string) error {
	if dbg {
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
	err = L.Call(0, 0)
	if err != nil {
		return err
	}
	return nil
}

func loadScripts(L lua.Lua) error {
	exeName, exeNameErr := dos.GetModuleFileName()
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
		fmt.Println("Copyright (c) 2014-2016 HAYAMA_Kaoru and NYAOS.ORG")
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
					if err := loadBundleScript1(L, relpath); err != nil {
						fmt.Fprintf(os.Stderr, "cached %s: %s\n", relpath, err)
					}
				} else {
					path1 := filepath.Join(nyagos_d, name1)
					if dbg {
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
			if err1 := loadBundleScript1(L, relpath); err1 != nil {
				fmt.Fprintf(os.Stderr, "bundled %s: %s\n", relpath, err1.Error())
			}
		}
	}
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	dot_nyagos := filepath.Join(home, ".nyagos")
	dotStat, dotErr := os.Stat(dot_nyagos)
	if dotErr != nil {
		return nil
	}
	cachePath := filepath.Join(AppDataDir(), "dotnyagos.luac")
	cacheStat, cacheErr := os.Stat(cachePath)
	if cacheErr == nil && !dotStat.ModTime().After(cacheStat.ModTime()) {
		if dbg {
			println("load cached ", cachePath)
		}
		if _, err := L.LoadFile(cachePath, "b"); err == nil {
			L.Call(0, 0)
			return nil
		}
	}
	if dbg {
		println("load real ", dot_nyagos)
	}
	if _, err := L.LoadFile(dot_nyagos, "bt"); err != nil {
		return err
	}
	chank := L.Dump()
	if err := L.Call(0, 0); err != nil {
		return err
	}
	w, w_err := os.Create(cachePath)
	if w_err != nil {
		return w_err
	}
	w.Write(chank)
	w.Close()
	return nil
}
