// +build ignore

package main

import (
	"fmt"

	"github.com/zetamatta/nyagos/dos"
)

func main() {
	netdrive, err := dos.GetNetDrives()
	if err != nil {
		panic(err.Error())
	}
	for _, d := range netdrive {
		fmt.Printf("net use %c: \"%s\"\n", d.Letter, d.Remote)
	}

	d, err := dos.FindVacantDrive()
	if err != nil {
		println(err)
		return
	}
	fmt.Printf("last drive=%c\n", d)
}
