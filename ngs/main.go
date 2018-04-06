package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/nyagos/mains"
)

func main() {
	mains.Version = "without Lua"

	if err := mains.Start(mains.Main); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
