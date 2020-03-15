// +build run

package main

import (
	"github.com/zetamatta/nyagos/nodos"
	"time"
)

func main() {
	d, err := nodos.OsDateLayout()
	if err != nil {
		println(err.Error())
		return
	}
	println("layout=", d)
	println("today=", time.Now().Format(d))
}
