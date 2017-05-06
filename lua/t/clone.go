package main

import (
	".."
)

func main() {
	L, err := lua.New()
	if err != nil {
		println(err)
		return
	}
	defer L.Close()
	L.OpenLibs()

	L.PushString("ahaha")
	L.SetGlobal("hogehoge")

	L.NewTable()
	L.PushInteger(1)
	L.SetField(-2, "foo")
	L.SetGlobal("sample_table")

	println("Setup first instance: done")

	n, err := L.Clone()
	if err != nil {
		println(err)
		return
	}
	N, ok := n.(lua.Lua)
	if !ok {
		println("Cast fail interface{} to lua.Lua")
		return
	}
	println("Cloning done")
	defer N.Close()
	N.GetGlobal("hogehoge")
	val, err := N.ToString(-1)
	if err != nil {
		println(err)
		return
	}
	L.Pop(1)
	if val == "ahaha" {
		println("test-1: ok")
	} else {
		println("test-1: ng: val=", val)
	}

	N.GetGlobal("sample_table")
	N.GetField(-1, "foo")
	val2, err := N.ToInteger(-1)
	L.Pop(2)
	if val2 == 1 {
		println("test-2: ok")
	} else {
		println("test-2: ng")
	}
}
