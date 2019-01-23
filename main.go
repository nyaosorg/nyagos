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
	rc := 0
	if err := frame.Start(mains.Main); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err)
		rc = 1
	}
	if defined.DBG {
		os.Stdin.Read(make([]byte, 1))
	}
	os.Exit(rc)
}
