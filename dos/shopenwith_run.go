// +build run

package main

import (
	"os"

	"github.com/zetamatta/nyagos/dos"
)

func main() {
	for _, arg := range os.Args[1:] {
		err := dos.ShOpenWithDialog(arg, "")
		if err != nil {
			println(err.Error())
		}
	}
}
