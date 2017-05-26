package main

import (
	"unsafe"

	"github.com/zetamatta/nyagos/lua"
)

func main() {
	src := []byte{'A', 'B', 'C', '\000'}
	s := lua.CGoStringZ(uintptr(unsafe.Pointer(&src[0])))
	println(s)
	if s == "ABC" {
		println("->OK")
	} else {
		println("->NG")
	}
}
