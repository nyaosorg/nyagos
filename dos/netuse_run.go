// +build run

package main

import (
	"github.com/zetamatta/nyagos/dos"
)

func main() {
	err := dos.WNetAddConnection2(`\\localhost\C$`, "O:", "", "")
	if err != nil {
		println(err.Error())
	}
}
