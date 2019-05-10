package main

import (
	"fmt"
	"time"

	"github.com/zetamatta/nyagos/dos"
)

func main() {
	indent := 1
	var callback func(*dos.NetResource) bool

	callback = func(node *dos.NetResource) bool {
		name := node.RemoteName()
		now := time.Now()
		fmt.Printf("%02d:%02d:%02d %*s%s\n",
			now.Hour(),
			now.Minute(),
			now.Second(),
			indent*2,
			"",
			name)
		if len(name) <= 0 || name[0] != '\\' {
			indent++
			node.Enum(callback)
			indent--
		}
		return true
	}
	err := dos.WNetEnum(callback)
	if err != nil {
		println(err.Error())
	}
}
