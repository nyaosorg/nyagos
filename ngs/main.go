package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/nyagos/frame"
)

func main() {
	frame.Version = "without Lua"

	if err := frame.Start(Main); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
