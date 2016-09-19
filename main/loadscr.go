package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

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

func loadScripts(L lua.Lua) {
	exeName, exeNameErr := dos.GetModuleFileName()
	if exeNameErr != nil {
		fmt.Fprintln(os.Stderr, exeNameErr)
	}
	exeFolder := filepath.Dir(exeName)

	// for compatibility
	nyagos_lua := filepath.Join(exeFolder, "nyagos.lua")
	if _, err := os.Stat(nyagos_lua); err == nil {
		err := L.Source(nyagos_lua)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		return
	}

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
		names, err := nyagos_d_fd.Readdirnames(-1)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		} else {
			for _, name1 := range names {
				if strings.HasSuffix(strings.ToLower(name1), ".lua") {
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
