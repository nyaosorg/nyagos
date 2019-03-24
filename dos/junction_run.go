// +build run

package main

import (
	"os"

	"github.com/zetamatta/nyagos/dos"
)

func main() {
	if len(os.Args) < 3 {
		println("go run junction.go DST SRC")
		return
	}
	if err := dos.CreateJunction(os.Args[1], os.Args[2]); err != nil {
		println(err.Error())
		os.Exit(1)
	}
}
