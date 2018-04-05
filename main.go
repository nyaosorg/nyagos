package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/mainl"
	"github.com/zetamatta/nyagos/mains"
)

var stamp string
var commit string
var version string

func main() {
	mains.Version = version

	if err := mains.Start(mainl.Main); err != nil {
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
