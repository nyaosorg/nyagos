package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/mains"
)

func when_panic() {
	err := recover()
	if err == nil {
		return
	}
	var dump bytes.Buffer
	w := io.MultiWriter(os.Stderr, &dump)

	fmt.Fprintln(w, "************ Panic Occurred. ***********")
	fmt.Fprintln(w, err)
	w.Write(debug.Stack())
	fmt.Fprintln(w, "*** Please copy these error message ***")
	fmt.Fprintln(w, "*** And hit ENTER key to quit.      ***")

	ioutil.WriteFile("nyagos.dump", dump.Bytes(), 0666)

	var dummy [1]byte
	os.Stdin.Read(dummy[:])
}

var stamp string
var commit string
var version string

func main() {
	defer when_panic()

	if len(os.Args) >= 2 && os.Args[1] == "--show-version-only" {
		fmt.Printf("%s-%s\n", version, runtime.GOARCH)
		os.Exit(0)
	}

	mains.Stamp = stamp
	mains.Commit = commit
	mains.Version = version

	if err := mains.Main(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if err != io.EOF {
			if defined.DBG {
				var dummy [1]byte
				os.Stdin.Read(dummy[:])
			}
			os.Exit(1)
		}
	}
	if defined.DBG {
		var dummy [1]byte
		os.Stdin.Read(dummy[:])
	}
	os.Exit(0)
}
