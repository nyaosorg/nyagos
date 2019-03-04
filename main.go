package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/nyagos/defined"
	"github.com/zetamatta/nyagos/frame"
	"github.com/zetamatta/nyagos/mains"
)

var version string

func main() {
	frame.Version = version
	if err := frame.Start(mains.Main); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err)
		defer os.Exit(1)
	}
	if defined.DBG {
		os.Stdin.Read(make([]byte, 1))
	}
}
