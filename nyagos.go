package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/zetamatta/nyagos/mains"
)

func when_panic() {
	err := recover()
	if err == nil {
		return
	}
	var dump bytes.Buffer
	w := io.MultiWriter(os.Stderr, &dump)

	fmt.Fprintln(w, "************ Panic Occured. ***********")
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

	mains.Stamp = stamp
	mains.Commit = commit
	mains.Version = version

	if err := mains.Main(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if err != io.EOF {
			os.Exit(1)
		}
	}
	os.Exit(0)
}
