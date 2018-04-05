package main

import (
	"fmt"
	"io"
	"os"

	"github.com/zetamatta/nyagos/mains"
)

func main() {
	mains.Version = "without Lua"

	if err := mains.Start(mains.Main); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		if err != io.EOF {
			os.Exit(1)
		}
	}
	os.Exit(0)
}
