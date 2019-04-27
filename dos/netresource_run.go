// +build run

package main

import (
	"github.com/zetamatta/nyagos/dos"
)

func main() {
	machines := []string{}

	err := dos.EachMachine(func(node *dos.NetResource) bool {
		machines = append(machines, node.RemoteName())
		return true
	})
	if err != nil {
		println(err.Error())
	}

	for _, name := range machines {
		println("machine:", name)
		err = dos.EachMachineNode(name, func(node *dos.NetResource) bool {
			println("  ", node.RemoteName())
			return true
		})
	}

	if err != nil {
		println(err.Error())
	}
}

// https://msdn.microsoft.com/ja-jp/library/cc447030.aspx
// http://eternalwindows.jp/security/share/share06.html
