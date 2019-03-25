// +build run

package main

import (
	"github.com/zetamatta/nyagos/dos"
)

func main() {
	label, fsname, err := dos.VolumeName("C:\\")
	if err != nil {
		println(err.Error())
		return
	}
	println(label)
	println(fsname)
}
