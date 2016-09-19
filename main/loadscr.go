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
	// print("try bundle version: ", path, "\n")
	bin, err := Asset(path)
	if err != nil {
		return err
	}
	err = L.LoadString(string(bin))
	if err != nil {
		return err
	}
	err = L.Call(0, 0)
	if err != nil {
		return err
	}
	return nil
}

func loadScripts(L lua.Lua) {
	exeName, exeNameErr := dos.GetModuleFileName()
	if exeNameErr != nil {
		fmt.Fprintln(os.Stderr, exeNameErr)
	}
	exeFolder := filepath.Dir(exeName)

	if !silentmode {
		fmt.Printf("Nihongo Yet Another GOing Shell %s-%s Powered by %s & %s\n",
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
				if assetErr == nil && !asset1.ModTime().Truncate(time.Second).Before(finfo1.ModTime().Truncate(time.Second)) {
					if err := loadBundleScript1(L, relpath); err != nil {
						fmt.Fprintf(os.Stderr, "cached %s: %s\n", relpath, err)
					}
				} else {
					path1 := filepath.Join(nyagos_d, name1)
					err1 := L.Source(path1)
					if err1 != nil {
						fmt.Fprintf(os.Stderr, "%s: %s\n", name1, err1.Error())
					}
				}
			}
		}
	}
	home := os.Getenv("HOME")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	dot_nyagos := filepath.Join(home, ".nyagos")
	if _, err := os.Stat(dot_nyagos); err == nil {
		err1 := L.Source(dot_nyagos)
		if err1 != nil {
			fmt.Fprintf(os.Stderr, "%s: %s\n", dot_nyagos, err1.Error())
		}
	}
}
