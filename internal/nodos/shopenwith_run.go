//go:build run
// +build run

package main

import (
	"os"

	"github.com/nyaosorg/nyagos/internal/nodos"
)

func main() {
	for _, arg := range os.Args[1:] {
		err := nodos.ShOpenWithDialog(arg, "")
		if err != nil {
			println(err.Error())
		}
	}
}
