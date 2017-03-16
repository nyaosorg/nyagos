package main

import (
	"fmt"
	"os"
	"runtime/debug"
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

func main() {
	defer when_panic()

	if err := Main(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
