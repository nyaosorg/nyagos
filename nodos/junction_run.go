// +build run

package main

import (
	"os"

	"github.com/zetamatta/nyagos/nodos"
)

func main() {
	if len(os.Args) < 3 {
		println("go run junction.go SRC DST")
		return
	}
	if err := nodos.CreateJunction(os.Args[1], os.Args[2]); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
