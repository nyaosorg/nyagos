package main

import (
	"fmt"
	"io"
	"os"

	"github.com/yuin/gopher-lua"
	"github.com/zetamatta/nyagos/frame"
)

func main() {
	frame.Version = "with " + lua.PackageName + "-" + lua.PackageVersion

	if err := frame.Start(Main); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
