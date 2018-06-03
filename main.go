package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/frame"
	mains "github.com/zetamatta/nyagos/gopherSh"
)

var version string

func main() {
	var dummy [1]byte
	frame.Version = version

	if err := frame.Start(mains.Main); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err.Error())
		if defined.DBG {
			os.Stdin.Read(dummy[:])
		}
		os.Exit(1)
	}
	if defined.DBG {
		os.Stdin.Read(dummy[:])
	}
	os.Exit(0)
}
