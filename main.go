package main

import (
	"embed"
	"fmt"
	"io"
	"os"

	"github.com/nyaosorg/nyagos/internal/defined"
	"github.com/nyaosorg/nyagos/internal/frame"
	"github.com/nyaosorg/nyagos/internal/mains"
	"github.com/nyaosorg/nyagos/internal/onexit"
)

var version = "snapshot"

//go:embed embed/*.lua
var embedLua embed.FS

func run() error {
	defer frame.PanicHandler()
	defer onexit.Done()
	if defined.DBG {
		defer os.Stdin.Read(make([]byte, 1))
	}
	frame.Setup(version)
	return mains.Run(&embedLua)
}

func main() {
	if err := run(); err != nil && err != io.EOF {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
