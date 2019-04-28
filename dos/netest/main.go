package main

import (
	"fmt"
	"github.com/zetamatta/nyagos/dos"
)

func main() {
	indent := 1
	var callback func(*dos.NetResource) bool

	callback = func(node *dos.NetResource) bool {
		fmt.Printf("%*s%s\n", indent*2, "", node.RemoteName())
		indent++
		node.Enum(callback)
		indent--
		return true
	}
	err := dos.WNetEnum(callback)
	if err != nil {
		println(err.Error())
	}
}
