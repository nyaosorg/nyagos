package main

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"./mains"
)

func when_panic() {
	err := recover()
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, "************ Panic Occured. ***********")
	fmt.Fprintln(os.Stderr, err)
	debug.PrintStack()
	fmt.Fprintln(os.Stderr, "*** Please copy these error message ***")
	fmt.Fprintln(os.Stderr, "*** And hit ENTER key to quit.      ***")
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
